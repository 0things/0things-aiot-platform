package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"0things/pkg/event"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestEventService_HandleEvent(t *testing.T) {
	logger := zap.NewNop()
	svc := NewEventService(viper.New(), logger)

	msg := event.DeviceMessage{
		DeviceKey:   "dev_event_01",
		Transport:   "mqtt",
		MessageType: "event",
		Payload:     json.RawMessage(`{"event_type": "online", "ip": "192.168.1.10"}`),
		Timestamp:   time.Now(),
	}

	if err := svc.HandleEvent(context.Background(), msg); err != nil {
		t.Errorf("HandleEvent returned error: %v", err)
	}
}
