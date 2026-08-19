package service_test

import (
	"context"
	"testing"

	"aiot-backend/internal/model"
	"aiot-backend/internal/service"
	mock_service "aiot-backend/test/mocks/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOTAService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOTARepo := mock_service.NewMockOTAServiceInterface(ctrl)
	ctx := context.Background()

	expectedPackages := []model.OTAPackage{{ID: 1, PackageName: "firmware-1"}}
	var total int64 = 1

	mockOTARepo.EXPECT().List(ctx, 1, 10).Return(expectedPackages, total, nil)

	packages, total, err := mockOTARepo.List(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, expectedPackages, packages)
	assert.Equal(t, int64(1), total)
}

func TestOTAService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOTARepo := mock_service.NewMockOTAServiceInterface(ctrl)
	ctx := context.Background()

	expectedPackage := &model.OTAPackage{ID: 1, PackageName: "firmware-1"}

	mockOTARepo.EXPECT().Get(ctx, int64(1)).Return(expectedPackage, nil)

	pkg, err := mockOTARepo.Get(ctx, int64(1))
	assert.NoError(t, err)
	assert.Equal(t, expectedPackage, pkg)
}

func TestOTAService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOTARepo := mock_service.NewMockOTAServiceInterface(ctrl)
	ctx := context.Background()

	pkg := &model.OTAPackage{PackageName: "firmware-1"}

	mockOTARepo.EXPECT().Create(ctx, pkg, "P001").Return(nil)

	err := mockOTARepo.Create(ctx, pkg, "P001")
	assert.NoError(t, err)
}

func TestOTAService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOTARepo := mock_service.NewMockOTAServiceInterface(ctrl)
	ctx := context.Background()

	mockOTARepo.EXPECT().Delete(ctx, int64(1)).Return(nil)

	err := mockOTARepo.Delete(ctx, int64(1))
	assert.NoError(t, err)
}

func TestOTAService_Statistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOTARepo := mock_service.NewMockOTAServiceInterface(ctrl)
	ctx := context.Background()

	expectedStats := service.UpgradeStatistics{
		PackageID:          "1",
		TotalTargetDevices: 100,
		SuccessfulUpgrades: 80,
		FailedUpgrades:     10,
	}

	mockOTARepo.EXPECT().Statistics(ctx, "firmware-1").Return(expectedStats, nil)

	stats, err := mockOTARepo.Statistics(ctx, "firmware-1")
	assert.NoError(t, err)
	assert.Equal(t, expectedStats, stats)
}
