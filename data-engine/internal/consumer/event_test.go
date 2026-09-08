package consumer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"0things/pkg/event"
	"data-engine/internal/service"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestEventConsumer_HandleEvent(t *testing.T) {
	v := viper.New()
	logger := zap.NewNop()
	eventService := service.NewEventService(v, logger)
	consumer := NewEventConsumer(eventService, logger)

	rawPayload := []byte(`{"alarm": "overvoltage", "val": 380}`)
	msg := &event.DeviceMessage{
		DeviceKey:   "dev_event_01",
		ProductKey:  "prod_01",
		Transport:   "mqtt",
		MessageType: "event",
		Payload:     json.RawMessage(rawPayload),
		Timestamp:   time.Now().UTC(),
	}

	err := consumer.HandleEvent(context.Background(), msg, nil)
	if err != nil {
		t.Fatalf("unexpected error in HandleEvent: %v", err)
	}
}
