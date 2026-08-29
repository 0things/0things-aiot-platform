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

	repo := NewDeviceGroupRepository(&IoTDB{DB: db})
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

	repo := NewDeviceGroupRepository(&IoTDB{DB: db})
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
}
