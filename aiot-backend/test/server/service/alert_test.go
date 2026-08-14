package service_test

import (
	"context"
	"testing"

	"0things-backend/internal/model"
	mock_service "0things-backend/test/mocks/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAlertService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlertRepo := mock_service.NewMockAlertServiceInterface(ctrl)
	ctx := context.Background()

	expectedAlerts := []model.Alert{{ID: 1, Summary: "High Temperature"}}
	var total int64 = 1

	mockAlertRepo.EXPECT().List(ctx, 1, 10, "", "", "").Return(expectedAlerts, total, nil)

	alerts, total, err := mockAlertRepo.List(ctx, 1, 10, "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, expectedAlerts, alerts)
	assert.Equal(t, int64(1), total)
}

func TestAlertService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlertRepo := mock_service.NewMockAlertServiceInterface(ctrl)
	ctx := context.Background()

	expectedAlert := &model.Alert{ID: 1, Summary: "High Temperature"}

	mockAlertRepo.EXPECT().Get(ctx, int64(1)).Return(expectedAlert, nil)

	alert, err := mockAlertRepo.Get(ctx, int64(1))
	assert.NoError(t, err)
	assert.Equal(t, expectedAlert, alert)
}

func TestAlertService_SetStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlertRepo := mock_service.NewMockAlertServiceInterface(ctrl)
	ctx := context.Background()

	expectedAlert := &model.Alert{ID: 1, Summary: "High Temperature", Status: "acknowledged"}

	mockAlertRepo.EXPECT().SetStatus(ctx, int64(1), "acknowledged").Return(expectedAlert, nil)

	alert, err := mockAlertRepo.SetStatus(ctx, int64(1), "acknowledged")
	assert.NoError(t, err)
	assert.Equal(t, expectedAlert, alert)
}
