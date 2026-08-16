package repository

import (
	"context"
	"testing"

	"0things-backend/internal/model"
	"0things-backend/internal/tenant"
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

func TestDeviceRepositoryListScopesSearchToTenant(t *testing.T) {
	store := newRepositoryTestDB(t, &model.Product{}, &model.Device{}, &model.DeviceState{})
	repo := &DeviceRepository{db: store.DB}
	ctx := context.Background()

	product1 := &model.Product{ProductKey: "P001", TenantID: 1}
	product2 := &model.Product{ProductKey: "P002", TenantID: 2}
	require.NoError(t, store.Create(product1).Error)
	require.NoError(t, store.Create(product2).Error)
	device1 := &model.Device{DeviceKey: "D001", Name: "Tenant one", ProductID: product1.ID, TenantID: 1}
	device2 := &model.Device{DeviceKey: "D002", Name: "Tenant two", ProductID: product2.ID, TenantID: 2}
	require.NoError(t, store.Create(device1).Error)
	require.NoError(t, store.Create(device2).Error)
	require.NoError(t, store.Create(&model.DeviceState{DeviceKey: device1.DeviceKey, State: "online"}).Error)
	require.NoError(t, store.Create(&model.DeviceState{DeviceKey: device2.DeviceKey, State: "online"}).Error)

	items, total, err := repo.List(ctx, 1, 20, 0, nil, nil, "Tenant")
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)

	items, total, err = repo.List(tenant.WithTenant(ctx, 2), 1, 20, 0, nil, nil, "Tenant")
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)
}
