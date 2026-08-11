package data

import (
	"context"
	"encoding/json"
	"fmt"

	"telemetry-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type telemetryEventRepo struct {
	data   *Data
	logger *log.Helper
}

// NewTelemetryEventRepo creates a new telemetry event repository.
func NewTelemetryEventRepo(data *Data, logger log.Logger) biz.TelemetryEventRepo {
	return &telemetryEventRepo{
		data:   data,
		logger: log.NewHelper(logger),
	}
}

// Save stores a telemetry event in Redis.
// The latest event for each device is stored with key: telemetry:device:{device_key}:latest
func (r *telemetryEventRepo) Save(ctx context.Context, event *biz.TelemetryEvent) error {
	r.logger.Infof("Saving telemetry event: type=%s, product_key=%s, device_key=%s",
		event.Type, event.ProductKey, event.DeviceKey)

	// Validate device_key
	if event.DeviceKey == "" {
		r.logger.Error("device_key is required")
		return fmt.Errorf("device_key is required")
	}

	// Marshal event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		r.logger.Errorf("Failed to marshal event to JSON: %v", err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Store latest event for the device
	key := fmt.Sprintf("telemetry:device:%s:latest", event.DeviceKey)
	if err := r.data.rdb.Set(ctx, key, eventJSON, 0).Err(); err != nil {
		r.logger.Errorf("Failed to save event to Redis: %v", err)
		return fmt.Errorf("failed to save event to Redis: %w", err)
	}

	r.logger.Debugf("Successfully saved event to Redis with key: %s", key)
	return nil
}
