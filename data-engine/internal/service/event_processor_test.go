package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"data-engine/internal/model"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestEventProcessor_HandleEvent(t *testing.T) {
	logger := zap.NewNop()
	proc := NewEventProcessor(viper.New(), logger)

	msg := model.DeviceMessage{
		DeviceKey:   "dev_event_01",
		Transport:   "mqtt",
		MessageType: "event",
		Payload:     json.RawMessage(`{"event_type": "online", "ip": "192.168.1.10"}`),
		Timestamp:   time.Now(),
	}

	if err := proc.HandleEvent(context.Background(), msg); err != nil {
		t.Errorf("HandleEvent returned error: %v", err)
	}
}
