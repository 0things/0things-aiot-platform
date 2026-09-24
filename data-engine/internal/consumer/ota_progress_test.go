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

type mockDeviceUpgradeStatusRepository struct {
	batchID   string
	deviceKey string
	progress  int32
	version   string
}

func (m *mockDeviceUpgradeStatusRepository) UpdateProgress(_ context.Context, batchID, deviceKey, _ string, progress int32) error {
	m.batchID = batchID
	m.deviceKey = deviceKey
	m.progress = progress
	return nil
}

func (m *mockDeviceUpgradeStatusRepository) UpdateDeviceInfo(_ context.Context, deviceKey, version string) error {
	m.deviceKey = deviceKey
	m.version = version
	return nil
}

func TestOTAProgressConsumer_HandleProgressReport(t *testing.T) {
	logger := zap.NewNop()
	store := &mockDeviceUpgradeStatusRepository{}
	otaService := service.NewOTAService(store, logger)
	consumer := NewOTAProgressConsumer(handler.NewOTAHandler(otaService, logger))

	payload, _ := json.Marshal(event.OTAProgressPayload{
		BatchID:   "b-01",
		Progress:  50,
		Timestamp: time.Now().UnixMilli(),
	})

	msg := &event.DeviceMessage{
		DeviceKey:   "dev_ota_01",
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeOTAProgress,
		Payload:     payload,
		Timestamp:   time.Now().UnixMilli(),
	}

	err := consumer.HandleProgressReport(context.Background(), msg, nil)
	if err != nil {
		t.Fatalf("unexpected error in HandleProgressReport: %v", err)
	}

	if store.progress != int32(50) {
		t.Errorf("unexpected progress recorded: progress=%d", store.progress)
	}
}
