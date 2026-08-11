package service

import (
	"context"
	"time"

	eventv1 "telemetry-service/api/event/v1"
	"telemetry-service/internal/data/kafka"

	"github.com/go-kratos/kratos/v2/log"
)

// EventService handles HTTP/gRPC requests for publishing device events.
type EventService struct {
	eventv1.UnimplementedEventServiceServer

	producer *kafka.Producer
	logger   *log.Helper
}

// NewEventService creates a new event service.
func NewEventService(producer *kafka.Producer, logger log.Logger) *EventService {
	return &EventService{
		producer: producer,
		logger:   log.NewHelper(log.With(logger, "module", "service/event")),
	}
}

// PublishEvent publishes a device event to Kafka.
func (s *EventService) PublishEvent(ctx context.Context, req *eventv1.PublishEventRequest) (*eventv1.PublishEventReply, error) {
	if s.producer == nil {
		return &eventv1.PublishEventReply{
			Code:    503,
			Message: "Kafka producer not configured",
		}, nil
	}

	// Validate required fields
	if req.Type == "" {
		return &eventv1.PublishEventReply{
			Code:    400,
			Message: "type is required",
		}, nil
	}
	if req.ProductKey == "" {
		return &eventv1.PublishEventReply{
			Code:    400,
			Message: "product_key is required",
		}, nil
	}
	if req.DeviceKey == "" {
		return &eventv1.PublishEventReply{
			Code:    400,
			Message: "device_key is required",
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
	if err := s.producer.SendJSON(ctx, kafka.TopicEvent, key, payload); err != nil {
		s.logger.Errorf("Failed to publish event: %v", err)
		return &eventv1.PublishEventReply{
			Code:    500,
			Message: "Failed to publish message",
		}, nil
	}

	s.logger.Infof("Published event: type=%s, product_key=%s, device_key=%s",
		req.Type, req.ProductKey, req.DeviceKey)
	return &eventv1.PublishEventReply{
		Code:    200,
		Message: "Event published successfully",
	}, nil
}
