package consumer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"data-engine/internal/event"
	"data-engine/internal/tsdb"
	"data-engine/internal/handler"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

func TestTelemetryConsumer_HandlePropertyPost(t *testing.T) {
	logger := zap.NewNop()
	mockTSDB := tsdb.NewMockClient(logger)
	telemetryService := service.NewTelemetryService(logger, mockTSDB)

	consumer := NewTelemetryConsumer(handler.NewTelemetryHandler(telemetryService, logger))

	rawPayload, _ := json.Marshal([]event.DevicePropertyPostPayload{
		{
			Timestamp: time.Now().UnixMilli(),
			Values:    map[string]interface{}{"temperature": 26.5, "humidity": 65},
		},
	})
	msg := &event.DeviceMessage{
		DeviceKey:   "dev_telemetry_01",
		ProductKey:  "prod_01",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeTelemetry,
		Payload:     json.RawMessage(rawPayload),
		Timestamp:   time.Now().UnixMilli(),
	}

	err := consumer.HandlePropertyPost(context.Background(), msg, nil)
	if err != nil {
		t.Fatalf("unexpected error in HandlePropertyPost: %v", err)
	}
}
