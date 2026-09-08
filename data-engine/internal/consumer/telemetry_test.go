package consumer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"0things/pkg/event"
	"data-engine/internal/engine"
	"data-engine/internal/storage"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestTelemetryConsumer_HandleTelemetry(t *testing.T) {
	v := viper.New()
	logger := zap.NewNop()
	shadow := storage.NewShadowStore(v, logger)
	processor := engine.NewProcessor(v, logger, nil, shadow)

	consumer := NewTelemetryConsumer(processor, logger)

	rawPayload := []byte(`{"temperature": 26.5, "humidity": 65}`)
	msg := &event.DeviceMessage{
		DeviceKey:   "dev_telemetry_01",
		ProductKey:  "prod_01",
		Transport:   "mqtt",
		MessageType: "telemetry",
		Payload:     json.RawMessage(rawPayload),
		Timestamp:   time.Now().UTC(),
	}

	err := consumer.HandleTelemetry(context.Background(), msg, nil)
	if err != nil {
		t.Fatalf("unexpected error in HandleTelemetry: %v", err)
	}

	// Verify shadow was updated
	sh, err := shadow.GetShadow(context.Background(), "dev_telemetry_01")
	if err != nil {
		t.Fatalf("failed to get shadow: %v", err)
	}
	if sh.Attributes["humidity"] != float64(65) {
		t.Errorf("expected humidity 65, got %v", sh.Attributes["humidity"])
	}
}

func TestTelemetryConsumer_HandleAttribute(t *testing.T) {
	v := viper.New()
	logger := zap.NewNop()
	shadow := storage.NewShadowStore(v, logger)
	processor := engine.NewProcessor(v, logger, nil, shadow)

	consumer := NewTelemetryConsumer(processor, logger)

	rawPayload := []byte(`{"ip": "192.168.1.100", "battery": 90}`)
	msg := &event.DeviceMessage{
		DeviceKey:   "dev_attr_01",
		ProductKey:  "prod_01",
		Transport:   "http",
		MessageType: "attributes",
		Payload:     json.RawMessage(rawPayload),
		Timestamp:   time.Now().UTC(),
	}

	err := consumer.HandleAttribute(context.Background(), msg, nil)
	if err != nil {
		t.Fatalf("unexpected error in HandleAttribute: %v", err)
	}

	sh, err := shadow.GetShadow(context.Background(), "dev_attr_01")
	if err != nil {
		t.Fatalf("failed to get shadow: %v", err)
	}
	if sh.Attributes["ip"] != "192.168.1.100" {
		t.Errorf("expected ip 192.168.1.100, got %v", sh.Attributes["ip"])
	}
}
