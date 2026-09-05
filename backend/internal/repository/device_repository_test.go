package repository

import (
	"context"
	"testing"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"

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

func TestDeviceRepository_TenantIsolation(t *testing.T) {
	store := newRepositoryTestDB(t, &model.Product{}, &model.Device{}, &model.DeviceState{})
	repo := NewDeviceRepository(store, nil)
	ctx := tenant.WithTenant(context.Background(), 1)

	product1 := &model.Product{ProductKey: "P001", Name: "Product one", OrganizationID: 1}
	product2 := &model.Product{ProductKey: "P002", Name: "Product two", OrganizationID: 2}
	require.NoError(t, store.Create(product1).Error)
	require.NoError(t, store.Create(product2).Error)
	device1 := &model.Device{DeviceKey: "D001", Name: "Tenant one", ProductID: product1.ID, OrganizationID: 1}
	device2 := &model.Device{DeviceKey: "D002", Name: "Tenant two", ProductID: product2.ID, OrganizationID: 2}
	require.NoError(t, store.Create(device1).Error)
	require.NoError(t, store.Create(device2).Error)
	require.NoError(t, store.Create(&model.DeviceState{DeviceKey: device1.DeviceKey, State: "online"}).Error)
	require.NoError(t, store.Create(&model.DeviceState{DeviceKey: device2.DeviceKey, State: "online"}).Error)

	items, total, err := repo.List(ctx, dto.ListDevicesQuery{Page: 1, PageSize: 20, Search: "Tenant"})
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)

	items, total, err = repo.List(tenant.WithTenant(ctx, 2), dto.ListDevicesQuery{Page: 1, PageSize: 20, Search: "Tenant"})
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)
}
