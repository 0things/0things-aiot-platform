package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"0things/pkg/event"
	"0things/pkg/protocol"
	"0things/pkg/tsdb"
	"data-engine/internal/model"
	"data-engine/internal/repository"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// TelemetryService processes device telemetry metrics, writes to TSDB, maintains device shadow state, and evaluates alarm rules.
type TelemetryService interface {
	ProcessMessage(ctx context.Context, msg event.DeviceMessage) error
}

type telemetryService struct {
	tsdbClient tsdb.Client
	shadowRepo repository.ShadowRepository
	protocols  *protocol.Registry
	logger     *zap.Logger
}

// NewTelemetryService creates an instance of TelemetryService.
func NewTelemetryService(config *viper.Viper, logger *zap.Logger, tsdbClient tsdb.Client, shadowRepo repository.ShadowRepository) TelemetryService {
	return &telemetryService{
		tsdbClient: tsdbClient,
		shadowRepo: shadowRepo,
		protocols:  protocol.DefaultRegistry(),
		logger:     logger,
	}
}

// ProcessMessage decodes payload, extracts metric records, persists to TSDB, updates shadow, and triggers rules.
func (s *telemetryService) ProcessMessage(ctx context.Context, msg event.DeviceMessage) error {
	// 1. Decode payload via application codec or fallback JSON unmarshal
	var data map[string]interface{}
	codec, ok := s.protocols.Get("json")
	if ok {
		decoded, err := codec.Decode(ctx, msg.Payload)
		if err == nil {
			data = decoded
		}
	}

	if data == nil {
		if jsonErr := json.Unmarshal(msg.Payload, &data); jsonErr != nil {
			s.logger.Warn("payload decoding yielded empty map, raw bytes skipped", zap.String("device_key", msg.DeviceKey))
			return nil
		}
	}

	// Unwrap nested properties/params/values map if present
	if params, ok := data["params"].(map[string]interface{}); ok {
		data = params
	} else if values, ok := data["values"].(map[string]interface{}); ok {
		data = values
	}

	// 2. Extract standard metric records
	records := make([]tsdb.Record, 0, len(data))
	ts := msg.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	for k, v := range data {
		if k == "params" || k == "values" || k == "properties" {
			if subMap, ok := v.(map[string]interface{}); ok {
				for subK, subV := range subMap {
					records = append(records, tsdb.Record{
						DeviceKey: msg.DeviceKey,
						Metric:    subK,
						Value:     subV,
						Timestamp: ts,
					})
					s.evaluateRule(msg.DeviceKey, subK, subV)
				}
				continue
			}
		}

		records = append(records, tsdb.Record{
			DeviceKey: msg.DeviceKey,
			Metric:    k,
			Value:     v,
			Timestamp: ts,
		})
		s.evaluateRule(msg.DeviceKey, k, v)
	}

	// 3. Batch write extracted metrics to pluggable TSDB client
	if len(records) > 0 && s.tsdbClient != nil {
		if err := s.tsdbClient.WriteBatch(ctx, records); err != nil {
			s.logger.Error("failed to write records to TSDB client", zap.String("device_key", msg.DeviceKey), zap.Error(err))
		}
	}

	// 4. Update device shadow snapshot
	if s.shadowRepo != nil && len(data) > 0 {
		_ = s.shadowRepo.UpdateShadow(ctx, msg.DeviceKey, data, msg.Timestamp)
	}

	s.logger.Debug("extracted & stored telemetry metrics via tsdb.Client", zap.String("device_key", msg.DeviceKey), zap.Int("count", len(records)))
	return nil
}

// evaluateRule evaluates metric thresholds and emits alarm events.
func (s *telemetryService) evaluateRule(deviceKey, metric string, val interface{}) {
	if metric == "temperature" || metric == "temp" {
		floatVal, ok := parseNumericValue(val)
		if !ok {
			return
		}

		if floatVal > 70.0 {
			alarm := model.AlarmEvent{
				DeviceKey:   deviceKey,
				RuleName:    "High Temperature Alarm",
				Level:       "CRITICAL",
				Description: fmt.Sprintf("device temperature reached %.1f°C exceeds threshold 70.0°C", floatVal),
				Timestamp:   time.Now().UTC(),
			}
			s.logger.Warn("🚨 RULE TRIGGERED: Alarm generated!",
				zap.String("device_key", alarm.DeviceKey),
				zap.String("rule", alarm.RuleName),
				zap.String("level", alarm.Level),
				zap.String("desc", alarm.Description),
			)
		}
	}
}

// parseNumericValue robustly converts arbitrary numeric data types to float64.
func parseNumericValue(val interface{}) (float64, bool) {
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case json.Number:
		if f, err := v.Float64(); err == nil {
			return f, true
		}
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}
