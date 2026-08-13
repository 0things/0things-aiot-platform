package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"0things-backend/internal/model"
	"github.com/stretchr/testify/require"
)

func TestDeviceEventRepositoryList(t *testing.T) {
	store := newRepositoryTestDB(t, &model.Product{}, &model.Device{}, &model.DeviceEvent{})
	ctx := context.Background()
	product := &model.Product{ProductKey: "sensor", Name: "Sensor"}
	require.NoError(t, store.Create(product).Error)
	device := &model.Device{DeviceKey: "D001", Name: "Workshop sensor", ProductID: product.ID}
	require.NoError(t, store.Create(device).Error)
	repo := NewDeviceEventRepository(store)
	older := time.Now().Add(-time.Hour)
	require.NoError(t, repo.Create(ctx, &model.DeviceEvent{DeviceID: device.ID, EventType: "lowBattery", EventAt: older, Data: json.RawMessage(`{"level":10}`)}))
	require.NoError(t, repo.Create(ctx, &model.DeviceEvent{DeviceID: device.ID, EventType: "overheat", EventAt: time.Now(), Data: json.RawMessage(`{"value":90}`)}))

	events, total, err := repo.List(ctx, 1, 20, "workshop", "D001", "", nil, nil)
	require.NoError(t, err)
	require.Equal(t, int64(2), total)
	require.Equal(t, "overheat", events[0].EventType)
	require.Equal(t, "D001", events[0].DeviceKey)

	startAt := time.Now().Add(-30 * time.Minute)
	events, total, err = repo.List(ctx, 1, 20, "", "", "over", &startAt, nil)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, events, 1)
}
