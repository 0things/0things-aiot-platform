package repository

import (
	"context"
	"testing"

	"0things-backend/internal/model"
	"github.com/stretchr/testify/require"
)

func TestDeviceRepositories(t *testing.T) {
	store := newRepositoryTestDB(t,
		&model.DeviceTag{},
		&model.DeviceShadow{},
		&model.DeviceShadowHistory{},
		&model.DevicePushRecord{},
	)
	tags := NewDeviceTagRepository(store)
	shadows := NewDeviceShadowRepository(store)
	pushRecords := NewPushRecordRepository(store)
	ctx := context.Background()

	require.NoError(t, tags.SetTags(ctx, 11, map[string]string{"region": "cn"}, true))
	storedTags, err := tags.ListTags(ctx, 11)
	require.NoError(t, err)
	require.Len(t, storedTags, 1)
	require.Equal(t, "region", storedTags[0].Key)

	desired := map[string]any{"power": true}
	shadow, err := shadows.MutateShadow(ctx, 11, 0, "app", &desired, nil, false)
	require.NoError(t, err)
	require.EqualValues(t, 1, shadow.Version)
	history, err := shadows.ListShadowHistory(ctx, 11)
	require.NoError(t, err)
	require.Len(t, history, 1)

	record := &model.DevicePushRecord{DeviceID: 11, Status: "success"}
	require.NoError(t, pushRecords.CreatePushRecord(ctx, record))
	storedRecords, total, err := pushRecords.ListPushRecords(ctx, 11, 1, 20, "", "")
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, storedRecords, 1)
}
