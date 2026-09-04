package service_test

import (
	"context"
	"testing"

	"aiot-backend/internal/dto"
	mock_service "aiot-backend/test/mocks/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeviceEventService_Record(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventRepo := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	ctx := context.Background()

	data := map[string]any{"temperature": 25.5}

	mockEventRepo.EXPECT().Record(ctx, "P001", "D001", "temperature", int64(0), data).Return(nil)

	err := mockEventRepo.Record(ctx, "P001", "D001", "temperature", int64(0), data)
	assert.NoError(t, err)
}

func TestDeviceEventService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventRepo := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	ctx := context.Background()

	expectedEvents := []dto.DeviceEventListItem{{ID: 1, EventType: "temperature"}}
	var total int64 = 1
	query := dto.ListDeviceEventsQuery{}
	mockEventRepo.EXPECT().List(ctx, query).Return(expectedEvents, total, nil)

	events, total, err := mockEventRepo.List(ctx, query)
	assert.NoError(t, err)
	assert.Equal(t, expectedEvents, events)
	assert.Equal(t, int64(1), total)
}
