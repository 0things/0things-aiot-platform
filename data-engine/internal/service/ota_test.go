package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"0things/pkg/event"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockDeviceUpgradeStatusRepository struct {
	batchID   string
	deviceKey string
	version   string
	errDesc   string
	progress  int32
	updateErr error
}

func (m *mockDeviceUpgradeStatusRepository) UpdateProgress(_ context.Context, batchID, deviceKey, errDesc string, progress int32) error {
	m.batchID = batchID
	m.deviceKey = deviceKey
	m.errDesc = errDesc
	m.progress = progress
	return m.updateErr
}

func (m *mockDeviceUpgradeStatusRepository) UpdateDeviceInfo(_ context.Context, deviceKey, version string) error {
	m.deviceKey = deviceKey
	m.version = version
	return m.updateErr
}

func TestOTAService_HandleProgressReport(t *testing.T) {
	ctx := context.Background()

	t.Run("valid progress report", func(t *testing.T) {
		repo := &mockDeviceUpgradeStatusRepository{}
		svc := NewOTAService(repo, zap.NewNop())

		payload, _ := json.Marshal(event.OTAProgressPayload{
			BatchID:   "b_01",
			Progress:  65,
			Timestamp: time.Now().UnixMilli(),
		})

		err := svc.HandleProgressReport(ctx, event.DeviceMessage{
			DeviceKey:   "dev_01",
			ProductKey:  "pk_01",
			Transport:   event.TransportMQTT,
			MessageType: event.MessageTypeOTAProgress,
			Payload:     payload,
			Timestamp:   time.Now().UnixMilli(),
		})
		require.NoError(t, err)
		assert.Equal(t, "b_01", repo.batchID)
		assert.Equal(t, "dev_01", repo.deviceKey)
		assert.Equal(t, int32(65), repo.progress)
	})

	t.Run("error report", func(t *testing.T) {
		repo := &mockDeviceUpgradeStatusRepository{}
		svc := NewOTAService(repo, zap.NewNop())

		payload, _ := json.Marshal(event.OTAProgressPayload{
			BatchID:      "b_03",
			ErrorMessage: "FLASH_ERR: disk full",
		})

		err := svc.HandleProgressReport(ctx, event.DeviceMessage{
			DeviceKey:   "dev_03",
			Transport:   event.TransportMQTT,
			MessageType: event.MessageTypeOTAProgress,
			Payload:     payload,
			Timestamp:   time.Now().UnixMilli(),
		})
		require.NoError(t, err)
		assert.Equal(t, "FLASH_ERR: disk full", repo.errDesc)
	})

	t.Run("missing required envelope fields", func(t *testing.T) {
		repo := &mockDeviceUpgradeStatusRepository{}
		svc := NewOTAService(repo, zap.NewNop())

		err := svc.HandleProgressReport(ctx, event.DeviceMessage{
			DeviceKey: "",
			Payload:   []byte(`{"batch_id":"b_04"}`),
			Timestamp: time.Now().UnixMilli(),
		})
		assert.NoError(t, err)
		assert.Empty(t, repo.batchID)
	})

	t.Run("missing required payload fields", func(t *testing.T) {
		repo := &mockDeviceUpgradeStatusRepository{}
		svc := NewOTAService(repo, zap.NewNop())

		err := svc.HandleProgressReport(ctx, event.DeviceMessage{
			DeviceKey:   "dev_04",
			Transport:   event.TransportMQTT,
			MessageType: event.MessageTypeOTAProgress,
			Payload:     []byte(`{"step": 10}`),
			Timestamp:   time.Now().UnixMilli(),
		})
		assert.NoError(t, err)
		assert.Empty(t, repo.batchID)
	})

	t.Run("invalid progress out of range", func(t *testing.T) {
		repo := &mockDeviceUpgradeStatusRepository{}
		svc := NewOTAService(repo, zap.NewNop())

		err := svc.HandleProgressReport(ctx, event.DeviceMessage{
			DeviceKey:   "dev_05",
			Transport:   event.TransportMQTT,
			MessageType: event.MessageTypeOTAProgress,
			Payload:     []byte(`{"batch_id":"b_05","step":150}`),
			Timestamp:   time.Now().UnixMilli(),
		})
		assert.NoError(t, err)
		assert.Empty(t, repo.batchID)
	})

	t.Run("repo error propagated", func(t *testing.T) {
		repo := &mockDeviceUpgradeStatusRepository{updateErr: errors.New("db error")}
		svc := NewOTAService(repo, zap.NewNop())

		payload, _ := json.Marshal(event.OTAProgressPayload{
			BatchID:  "b_06",
			Progress: 20,
		})

		err := svc.HandleProgressReport(ctx, event.DeviceMessage{
			DeviceKey:   "dev_06",
			Transport:   event.TransportMQTT,
			MessageType: event.MessageTypeOTAProgress,
			Payload:     payload,
			Timestamp:   time.Now().UnixMilli(),
		})
		assert.Error(t, err)
	})
}

func TestOTAService_HandleDeviceInfo(t *testing.T) {
	ctx := context.Background()

	t.Run("valid device info report", func(t *testing.T) {
		repo := &mockDeviceUpgradeStatusRepository{}
		svc := NewOTAService(repo, zap.NewNop())

		payload, _ := json.Marshal(event.OTAInformPayload{
			Version: "v2.0.0",
		})

		err := svc.HandleDeviceInfo(ctx, event.DeviceMessage{
			DeviceKey:   "dev_01",
			Transport:   event.TransportMQTT,
			MessageType: event.MessageTypeOTAInform,
			Payload:     payload,
			Timestamp:   time.Now().UnixMilli(),
		})
		require.NoError(t, err)
		assert.Equal(t, "dev_01", repo.deviceKey)
		assert.Equal(t, "v2.0.0", repo.version)
	})

	t.Run("missing required fields", func(t *testing.T) {
		repo := &mockDeviceUpgradeStatusRepository{}
		svc := NewOTAService(repo, zap.NewNop())

		err := svc.HandleDeviceInfo(ctx, event.DeviceMessage{
			DeviceKey:   "dev_01",
			Transport:   event.TransportMQTT,
			MessageType: event.MessageTypeOTAInform,
			Payload:     []byte(`{}`),
			Timestamp:   time.Now().UnixMilli(),
		})
		assert.NoError(t, err)
		assert.Empty(t, repo.deviceKey)
	})
}
