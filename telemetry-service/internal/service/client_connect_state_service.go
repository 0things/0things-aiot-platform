package service

import (
	"context"
	"time"

	clientConnectStatev1 "telemetry-service/api/client_connect_state/v1"
	"telemetry-service/internal/data/kafka"

	"github.com/go-kratos/kratos/v2/log"
)

// ClientConnectStateService handles HTTP/gRPC requests for publishing client connect state events.
type ClientConnectStateService struct {
	clientConnectStatev1.UnimplementedClientConnectStateServiceServer

	producer *kafka.Producer
	logger   *log.Helper
}

// NewClientConnectStateService creates a new client connect state service.
func NewClientConnectStateService(producer *kafka.Producer, logger log.Logger) *ClientConnectStateService {
	return &ClientConnectStateService{
		producer: producer,
		logger:   log.NewHelper(log.With(logger, "module", "service/client_connect_state")),
	}
}

// PublishClientConnectState publishes a client connect state event to Kafka.
func (s *ClientConnectStateService) PublishClientConnectState(ctx context.Context, req *clientConnectStatev1.PublishClientConnectStateRequest) (*clientConnectStatev1.PublishClientConnectStateReply, error) {
	if s.producer == nil {
		return &clientConnectStatev1.PublishClientConnectStateReply{
			Code:    503,
			Message: "Kafka producer not configured",
		}, nil
	}

	// Validate required fields
	if req.Clientid == "" {
		return &clientConnectStatev1.PublishClientConnectStateReply{
			Code:    400,
			Message: "clientid is required",
		}, nil
	}
	if req.EventType == "" {
		return &clientConnectStatev1.PublishClientConnectStateReply{
			Code:    400,
			Message: "event_type is required",
		}, nil
	}

	// Build payload
	timestamp := req.Timestamp
	if timestamp == 0 {
		timestamp = time.Now().UnixMilli()
	}
	payload := map[string]interface{}{
		"clientid":   req.Clientid,
		"timestamp":  timestamp,
		"event_type": req.EventType,
		"username":   req.Username,
		"reason":     req.Reason,
	}

	// Publish to Kafka
	if err := s.producer.SendJSON(ctx, kafka.TopicClientConnectState, req.Clientid, payload); err != nil {
		s.logger.Errorf("Failed to publish client connect state: %v", err)
		return &clientConnectStatev1.PublishClientConnectStateReply{
			Code:    500,
			Message: "Failed to publish message",
		}, nil
	}

	s.logger.Infof("Published client connect state event: clientid=%s, event_type=%s",
		req.Clientid, req.EventType)
	return &clientConnectStatev1.PublishClientConnectStateReply{
		Code:    200,
		Message: "Client connect state published successfully",
	}, nil
}
