package service_test

import (
	"context"
	"testing"

	"0things-backend/internal/model"
	mock_service "0things-backend/test/mocks/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRuleService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuleRepo := mock_service.NewMockRuleServiceInterface(ctrl)
	ctx := context.Background()

	expectedRules := []model.Rule{{ID: 1, Name: "Temperature Alert"}}
	var total int64 = 1

	mockRuleRepo.EXPECT().List(ctx, 1, 10, "", "", "").Return(expectedRules, total, nil)

	rules, total, err := mockRuleRepo.List(ctx, 1, 10, "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
	assert.Equal(t, int64(1), total)
}

func TestRuleService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuleRepo := mock_service.NewMockRuleServiceInterface(ctrl)
	ctx := context.Background()

	expectedRule := &model.Rule{ID: 1, Name: "Temperature Alert"}

	mockRuleRepo.EXPECT().Get(ctx, int64(1)).Return(expectedRule, nil)

	rule, err := mockRuleRepo.Get(ctx, int64(1))
	assert.NoError(t, err)
	assert.Equal(t, expectedRule, rule)
}

func TestRuleService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuleRepo := mock_service.NewMockRuleServiceInterface(ctrl)
	ctx := context.Background()

	rule := &model.Rule{Name: "Temperature Alert"}

	mockRuleRepo.EXPECT().Create(ctx, rule).Return(nil)

	err := mockRuleRepo.Create(ctx, rule)
	assert.NoError(t, err)
}

func TestRuleService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuleRepo := mock_service.NewMockRuleServiceInterface(ctrl)
	ctx := context.Background()

	mockRuleRepo.EXPECT().Delete(ctx, int64(1)).Return(nil)

	err := mockRuleRepo.Delete(ctx, int64(1))
	assert.NoError(t, err)
}

func TestRuleService_SetStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuleRepo := mock_service.NewMockRuleServiceInterface(ctrl)
	ctx := context.Background()

	expectedRule := &model.Rule{ID: 1, Name: "Temperature Alert", Status: "active"}

	mockRuleRepo.EXPECT().SetStatus(ctx, int64(1), "active").Return(expectedRule, nil)

	rule, err := mockRuleRepo.SetStatus(ctx, int64(1), "active")
	assert.NoError(t, err)
	assert.Equal(t, expectedRule, rule)
}
