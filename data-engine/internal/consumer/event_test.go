package consumer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"0things/pkg/event"
	"data-engine/internal/handler"
	"data-engine/internal/model"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

type mockDeviceEventRepository struct{}

func (m *mockDeviceEventRepository) Create(ctx context.Context, event *model.DeviceEvent) error {
	return nil
}

func TestEventConsumer_HandleEvent(t *testing.T) {
	logger := zap.NewNop()
	eventService := service.NewEventService(logger, &mockDeviceEventRepository{})
	consumer := NewEventConsumer(handler.NewEventHandler(eventService, logger))

	rawPayload := []byte(`{"identifier": "overvoltage", "type": "alert", "timestamp": 1726315200000, "data": {"val": 380}}`)
	msg := &event.DeviceMessage{
		DeviceKey:   "dev_event_01",
		ProductKey:  "prod_01",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeEvent,
		Payload:     json.RawMessage(rawPayload),
		Timestamp:   time.Now().UnixMilli(),
	}

	err := consumer.HandleEvent(context.Background(), msg, nil)
	if err != nil {
		t.Fatalf("unexpected error in HandleEvent: %v", err)
	}
}
