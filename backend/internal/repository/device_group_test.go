package repository

import (
	"context"
	"testing"

	"aiot-backend/internal/model"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDeviceGroupRepositoryDevices_DynamicRuleWithOR(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.Product{},
		&model.Device{},
		&model.DeviceState{},
		&model.DeviceGroup{},
		&model.DeviceGroupMember{},
	))

	products := []model.Product{
		{ProductKey: "product-a", Name: "A", OrganizationID: 1},
		{ProductKey: "product-b", Name: "B", OrganizationID: 1},
	}
	require.NoError(t, db.Create(&products).Error)
	require.NoError(t, db.Create(&[]model.Device{
		{DeviceKey: "device-a", Name: "A", ProductID: products[0].ID, OrganizationID: 1, Enabled: true},
		{DeviceKey: "device-b", Name: "B", ProductID: products[1].ID, OrganizationID: 1, Enabled: true},
	}).Error)

	repo := NewDeviceGroupRepository(db)
	group := &model.DeviceGroup{
		ID:             1,
		OrganizationID: 1,
		Type:           "dynamic",
		Rule:           "product_key = 'product-a' OR product_key = 'product-b'",
	}

	devices, total, err := repo.Devices(context.Background(), group)
	require.NoError(t, err)
	require.Equal(t, int64(2), total)
	require.Len(t, devices, 2)
}

func TestDeviceGroupRepositoryDevicesPage_Paginates(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Product{}, &model.Device{}, &model.DeviceState{}, &model.DeviceGroup{}, &model.DeviceGroupMember{}))
	product := model.Product{ProductKey: "product-a", Name: "A", OrganizationID: 1}
	require.NoError(t, db.Create(&product).Error)
	require.NoError(t, db.Create(&[]model.Device{
		{DeviceKey: "device-a", Name: "A", ProductID: product.ID, OrganizationID: 1, Enabled: true},
		{DeviceKey: "device-b", Name: "B", ProductID: product.ID, OrganizationID: 1, Enabled: true},
	}).Error)

	repo := NewDeviceGroupRepository(db)
	group := &model.DeviceGroup{Type: "dynamic", Rule: "product_key = 'product-a'"}
	devices, total, err := repo.DevicesPage(context.Background(), group, 2, 1, "", "")
	require.NoError(t, err)
	require.Equal(t, int64(2), total)
	require.Len(t, devices, 1)
}

func TestApplyGroupRule_RejectsUnsupportedField(t *testing.T) {
	_, err := applyGroupRule(nil, "password = 'secret'")
	require.EqualError(t, err, "unsupported group rule field: password")
	_, err = applyGroupRule(nil, "created_at = '2026-01-01'")
	require.EqualError(t, err, "unsupported group rule field: created_at")
	_, err = applyGroupRule(nil, `{"password":"secret"}`)
	require.EqualError(t, err, "unsupported group rule field: password")
}

func TestDeviceGroupRepositoryDevices_DynamicRuleWithJSON(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.Product{},
		&model.Device{},
		&model.DeviceState{},
		&model.DeviceTag{},
		&model.DeviceGroup{},
		&model.DeviceGroupMember{},
	))

	product := model.Product{ProductKey: "product-a", Name: "A", OrganizationID: 1}
	require.NoError(t, db.Create(&product).Error)
	devices := []model.Device{
		{ID: 1, DeviceKey: "device-1", Name: "D1", ProductID: product.ID, OrganizationID: 1, Enabled: true},
		{ID: 2, DeviceKey: "device-2", Name: "D2", ProductID: product.ID, OrganizationID: 1, Enabled: true},
	}
	require.NoError(t, db.Create(&devices).Error)
	require.NoError(t, db.Create(&[]model.DeviceState{
		{ID: 1, DeviceKey: "device-1", State: "online"},
		{ID: 2, DeviceKey: "device-2", State: "offline"},
	}).Error)
	require.NoError(t, db.Create(&[]model.DeviceTag{
		{DeviceID: 1, Key: "location", Value: "floor1"},
		{DeviceID: 2, Key: "location", Value: "floor2"},
	}).Error)

	repo := NewDeviceGroupRepository(db)

	// 1. JSON rule with state
	group := &model.DeviceGroup{
		ID:             1,
		OrganizationID: 1,
		Type:           "dynamic",
		Rule:           `{"state":"online"}`,
	}
	list, total, err := repo.Devices(context.Background(), group)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	require.Equal(t, "device-1", list[0].DeviceKey)

	// 2. JSON rule with tag
	groupTag := &model.DeviceGroup{
		ID:             2,
		OrganizationID: 1,
		Type:           "dynamic",
		Rule:           `{"tag.location":"floor2"}`,
	}
	listTag, totalTag, err := repo.Devices(context.Background(), groupTag)
	require.NoError(t, err)
	require.Equal(t, int64(1), totalTag)
	require.Len(t, listTag, 1)
	require.Equal(t, "device-2", listTag[0].DeviceKey)

	// 3. JSON rule with array
	groupArray := &model.DeviceGroup{
		ID:             3,
		OrganizationID: 1,
		Type:           "dynamic",
		Rule:           `{"state":["online","offline"]}`,
	}
	listArr, totalArr, err := repo.Devices(context.Background(), groupArray)
	require.NoError(t, err)
	require.Equal(t, int64(2), totalArr)
	require.Len(t, listArr, 2)
}

