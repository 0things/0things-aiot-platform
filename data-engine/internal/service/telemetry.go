package service

import (
	"context"
	"encoding/json"

	"data-engine/internal/event"
	"data-engine/internal/tsdb"

	"go.uber.org/zap"
)

// TelemetryService processes device telemetry/property metrics and writes to TSDB.
type TelemetryService interface {
	ProcessPropertyPost(ctx context.Context, msg event.DeviceMessage) error
}

type telemetryService struct {
	tsdbClient tsdb.Client
	logger     *zap.Logger
}

// NewTelemetryService creates an instance of TelemetryService.
func NewTelemetryService(logger *zap.Logger, tsdbClient tsdb.Client) TelemetryService {
	return &telemetryService{
		tsdbClient: tsdbClient,
		logger:     logger,
	}
}

// ProcessPropertyPost parses telemetry payload and writes metric records to TSDB.
func (s *telemetryService) ProcessPropertyPost(ctx context.Context, msg event.DeviceMessage) error {
	if len(msg.Payload) == 0 {
		return nil
	}

	var payloads []event.DevicePropertyPostPayload
	if err := json.Unmarshal(msg.Payload, &payloads); err != nil {
		s.logger.Warn("failed to unmarshal normalized telemetry payload",
			zap.String("device_key", msg.DeviceKey),
			zap.Error(err),
		)
		return nil
	}

	if len(payloads) == 0 {
		return nil
	}

	records := make([]tsdb.Record, 0)
	for _, p := range payloads {
		for metric, val := range p.Values {
			records = append(records, tsdb.Record{
				DeviceKey: msg.DeviceKey,
				Metric:    metric,
				Value:     val,
				Timestamp: p.Timestamp,
			})
		}
	}

	if len(records) > 0 {
		if err := s.tsdbClient.WriteBatch(ctx, records); err != nil {
			s.logger.Error("failed to write records to TSDB client", zap.String("device_key", msg.DeviceKey), zap.Error(err))
			return err
		}
	}

	s.logger.Debug("extracted & stored telemetry metrics via tsdb.Client", zap.String("device_key", msg.DeviceKey), zap.Int("count", len(records)))
	return nil
}
