package service

import (
	"context"
	"time"

	telemetryv1 "telemetry-service/api/telemetry/v1"
	"telemetry-service/internal/data/kafka"

	"github.com/go-kratos/kratos/v2/log"
)

// TelemetryService handles HTTP/gRPC requests for publishing telemetry events.
type TelemetryService struct {
	telemetryv1.UnimplementedTelemetryServiceServer

	producer *kafka.Producer
	logger   *log.Helper
}

// NewTelemetryService creates a new telemetry service.
func NewTelemetryService(producer *kafka.Producer, logger log.Logger) *TelemetryService {
	return &TelemetryService{
		producer: producer,
		logger:   log.NewHelper(log.With(logger, "module", "service/telemetry")),
	}
}

// PublishTelemetry publishes a telemetry event to Kafka.
func (s *TelemetryService) PublishTelemetry(ctx context.Context, req *telemetryv1.PublishTelemetryRequest) (*telemetryv1.PublishTelemetryReply, error) {
	if s.producer == nil {
		return &telemetryv1.PublishTelemetryReply{
			Code:    503,
			Message: "Kafka producer not configured",
		}, nil
	}

	// Validate required fields
	if req.Type == "" {
		return &telemetryv1.PublishTelemetryReply{
			Code:    400,
			Message: "type is required",
		}, nil
	}

	// Build payload
	timestamp := req.Timestamp
	if timestamp == 0 {
		timestamp = time.Now().UnixMilli()
	}
	payload := map[string]interface{}{
		"type":        req.Type,
		"timestamp":   timestamp,
		"product_key": req.ProductKey,
		"device_key":  req.DeviceKey,
		"session_id":  req.SessionId,
	}

	// Add data field if present
	if req.Data != nil {
		payload["data"] = req.Data.AsMap()
	}

	// Determine message key for partitioning
	key := req.ProductKey
	if key == "" && req.DeviceKey != "" {
		key = req.DeviceKey
	}

	// Publish to Kafka
	if err := s.producer.SendJSON(ctx, kafka.TopicTelemetry, key, payload); err != nil {
		s.logger.Errorf("Failed to publish telemetry: %v", err)
		return &telemetryv1.PublishTelemetryReply{
			Code:    500,
			Message: "Failed to publish message",
		}, nil
	}

	s.logger.Infof("Published telemetry event: type=%s, product_key=%s, device_key=%s",
		req.Type, req.ProductKey, req.DeviceKey)
	return &telemetryv1.PublishTelemetryReply{
		Code:    200,
		Message: "Telemetry published successfully",
	}, nil
}
