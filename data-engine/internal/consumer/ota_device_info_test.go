package consumer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"data-engine/internal/event"
	"data-engine/internal/handler"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

func TestOTADeviceInfoConsumer_HandleDeviceInfo(t *testing.T) {
	logger := zap.NewNop()
	store := &mockDeviceUpgradeStatusRepository{}
	otaService := service.NewOTAService(store, logger)
	consumer := NewOTADeviceInfoConsumer(handler.NewOTAHandler(otaService, logger))

	payload, _ := json.Marshal(event.OTAInformPayload{
		Version:   "v2.0.0",
		Timestamp: time.Now().UnixMilli(),
	})

	msg := &event.DeviceMessage{
		DeviceKey:   "dev_ota_01",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeOTAInform,
		Payload:     payload,
		Timestamp:   time.Now().UnixMilli(),
	}

	err := consumer.HandleDeviceInfo(context.Background(), msg, nil)
	if err != nil {
		t.Fatalf("unexpected error in HandleDeviceInfo: %v", err)
	}

	if store.version != "v2.0.0" {
		t.Errorf("unexpected version recorded: version=%s", store.version)
	}
}
