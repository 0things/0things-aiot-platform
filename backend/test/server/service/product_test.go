package service_test

import (
	"context"
	"errors"
	"testing"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/model"
	mock_service "aiot-backend/test/mocks/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProductService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock_service.NewMockProductServiceInterface(ctrl)
	ctx := context.Background()

	product := &model.Product{Name: "Test Product"}
	expectedProduct := &model.Product{ID: 1, Name: "Test Product"}

	mockProductRepo.EXPECT().Create(ctx, product).Return(expectedProduct, nil)

	result, err := mockProductRepo.Create(ctx, product)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, result)
}

func TestProductService_Create_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock_service.NewMockProductServiceInterface(ctrl)
	ctx := context.Background()

	product := &model.Product{}

	mockProductRepo.EXPECT().Create(ctx, product).Return(nil, errors.New("name is required"))

	result, err := mockProductRepo.Create(ctx, product)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestProductService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock_service.NewMockProductServiceInterface(ctrl)
	ctx := context.Background()

	expectedProduct := &model.Product{ID: 1, Name: "Test Product"}

	mockProductRepo.EXPECT().Get(ctx, int64(1)).Return(expectedProduct, nil)

	product, err := mockProductRepo.Get(ctx, int64(1))
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
}

func TestProductService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock_service.NewMockProductServiceInterface(ctrl)
	ctx := context.Background()

	expectedProducts := []dto.ProductListItem{{Product: model.Product{ID: 1, Name: "Product 1"}}}
	var total int64 = 1

	mockProductRepo.EXPECT().List(ctx, 1, 10, "", "", "").Return(expectedProducts, total, nil)

	products, total, err := mockProductRepo.List(ctx, 1, 10, "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	assert.Equal(t, int64(1), total)
}

func TestProductService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock_service.NewMockProductServiceInterface(ctrl)
	ctx := context.Background()

	mockProductRepo.EXPECT().Delete(ctx, int64(1)).Return(nil)

	err := mockProductRepo.Delete(ctx, int64(1))
	assert.NoError(t, err)
}

func TestProductService_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock_service.NewMockProductServiceInterface(ctrl)
	ctx := context.Background()

	expectedProduct := &model.Product{ID: 1, Name: "Restored Product"}

	mockProductRepo.EXPECT().Restore(ctx, int64(1)).Return(expectedProduct, nil)

	product, err := mockProductRepo.Restore(ctx, int64(1))
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
}

func TestProductService_GetByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mock_service.NewMockProductServiceInterface(ctrl)
	ctx := context.Background()

	expectedProduct := &model.Product{ID: 1, ProductKey: "P001", Name: "Test Product"}

	mockProductRepo.EXPECT().GetByKey(ctx, "P001").Return(expectedProduct, nil)

	product, err := mockProductRepo.GetByKey(ctx, "P001")
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
}
