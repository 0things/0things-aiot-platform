package biz

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
)

// ClientConnectStateEvent represents a client connect state event from EMQX.
type ClientConnectStateEvent struct {
	ClientID  string `json:"clientid"`
	Timestamp int64  `json:"timestamp"`
	Event     string `json:"event,omitempty"`      // client.connected, client.disconnected
	DeviceKey string `json:"device_key,omitempty"` // device identifier
	Reason    string `json:"reason,omitempty"`
}

// ClientConnectStateRepo defines the repository interface for storing client connect state events.
type ClientConnectStateRepo interface {
	Save(context.Context, *ClientConnectStateEvent) error
	UpdateDeviceOnline(ctx context.Context, deviceKey string, status string, timestamp int64) error
	UpdateDeviceOffline(ctx context.Context, deviceKey string, timestamp int64) error
}

// ClientConnectStateUsecase handles client connect state event processing.
type ClientConnectStateUsecase struct {
	repo   ClientConnectStateRepo
	logger *log.Helper
}

// NewClientConnectStateUsecase creates a new client connect state use case.
func NewClientConnectStateUsecase(repo ClientConnectStateRepo, logger log.Logger) *ClientConnectStateUsecase {
	return &ClientConnectStateUsecase{
		repo:   repo,
		logger: log.NewHelper(logger),
	}
}

// ProcessEvent processes a client connect state event (business logic).
func (uc *ClientConnectStateUsecase) ProcessEvent(ctx context.Context, event *ClientConnectStateEvent) error {
	// Business validation
	if event.ClientID == "" {
		return fmt.Errorf("clientid is required for client connect state events")
	}
	if event.Event == "" {
		return fmt.Errorf("event is required for client connect state events")
	}

	uc.logger.Infof("Processing client connect state event: clientid=%s, event=%s, device_key=%s",
		event.ClientID, event.Event, event.DeviceKey)

	// Process based on event type
	switch event.Event {
	case "client.connected":
		return uc.handleClientConnected(ctx, event)
	case "client.disconnected":
		return uc.handleClientDisconnected(ctx, event)
	default:
		uc.logger.Warnf("Unknown event type: %s", event.Event)
	}

	return nil
}

// handleClientConnected processes client connected events.
func (uc *ClientConnectStateUsecase) handleClientConnected(ctx context.Context, event *ClientConnectStateEvent) error {
	uc.logger.Infof("Client connected: clientid=%s, device_key=%s",
		event.ClientID, event.DeviceKey)

	// Update device status to online
	if event.DeviceKey != "" {
		if err := uc.repo.UpdateDeviceOnline(ctx, event.DeviceKey, "online", event.Timestamp); err != nil {
			uc.logger.Errorf("Failed to update device online status: %v", err)
			return err
		}
	}

	return nil
}

// handleClientDisconnected processes client disconnected events.
func (uc *ClientConnectStateUsecase) handleClientDisconnected(ctx context.Context, event *ClientConnectStateEvent) error {
	uc.logger.Infof("Client disconnected: clientid=%s, device_key=%s, reason=%s",
		event.ClientID, event.DeviceKey, event.Reason)

	// Update device status to offline
	if event.DeviceKey != "" {
		if err := uc.repo.UpdateDeviceOffline(ctx, event.DeviceKey, event.Timestamp); err != nil {
			uc.logger.Errorf("Failed to update device offline status: %v", err)
			return err
		}
	}

	return nil
}
