package service_test

import (
	"context"
	"testing"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"
	mock_service "aiot-backend/test/mocks/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeviceService_CreateDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mock_service.NewMockDeviceServiceInterface(ctrl)
	ctx := context.Background()

	device := &model.Device{Name: "Test Device", ProductID: 1}

	mockDeviceRepo.EXPECT().CreateDevice(ctx, device).Return(&model.Device{ID: 1, Name: "Test Device"}, nil)

	result, err := mockDeviceRepo.CreateDevice(ctx, device)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Device", result.Name)
}

func TestDeviceService_Device(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mock_service.NewMockDeviceServiceInterface(ctrl)
	ctx := context.Background()

	expectedDevice := &model.Device{ID: 1, Name: "Test Device"}

	mockDeviceRepo.EXPECT().Device(ctx, int64(1)).Return(expectedDevice, nil)

	device, err := mockDeviceRepo.Device(ctx, int64(1))
	assert.NoError(t, err)
	assert.Equal(t, expectedDevice, device)
}

func TestDeviceService_ListDevices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mock_service.NewMockDeviceServiceInterface(ctrl)
	ctx := context.Background()

	expectedDevices := []model.Device{{ID: 1, Name: "Device 1"}}
	var total int64 = 1
	query := dto.ListDevicesQuery{Page: 1, PageSize: 10}

	mockDeviceRepo.EXPECT().ListDevices(ctx, query).Return(expectedDevices, total, nil)

	devices, total, err := mockDeviceRepo.ListDevices(ctx, query)
	assert.NoError(t, err)
	assert.Equal(t, expectedDevices, devices)
	assert.Equal(t, int64(1), total)
}

func TestDeviceService_DeleteDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mock_service.NewMockDeviceServiceInterface(ctrl)
	ctx := context.Background()

	mockDeviceRepo.EXPECT().DeleteDevice(ctx, int64(1)).Return(nil)

	err := mockDeviceRepo.DeleteDevice(ctx, int64(1))
	assert.NoError(t, err)
}

func TestDeviceService_Stats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mock_service.NewMockDeviceServiceInterface(ctrl)
	ctx := context.Background()

	expectedStats := service.DeviceStatistics{
		TotalDevices:     10,
		ActivatedDevices: 8,
		OnlineDevices:    5,
		OfflineDevices:   3,
		InactiveDevices:  2,
	}

	mockDeviceRepo.EXPECT().Stats(ctx).Return(expectedStats, nil)

	stats, err := mockDeviceRepo.Stats(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedStats, stats)
}
