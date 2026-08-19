package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"aiot-backend/internal/handler"
	"aiot-backend/internal/model"
	mock_service "aiot-backend/test/mocks/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupProductRouter(mockService *mock_service.MockProductServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	productHandler := handler.NewProductHandler(h, mockService)

	router.POST("/products", productHandler.Create)
	router.GET("/products/:id", productHandler.Get)
	router.GET("/products/key/:productKey", productHandler.GetByKey)
	router.GET("/products", productHandler.List)
	router.PUT("/products/key/:productKey", productHandler.Update)
	router.DELETE("/products/:id", productHandler.Delete)
	router.POST("/products/:id/restore", productHandler.Restore)

	return router
}

func TestProductHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	router := setupProductRouter(mockService)

	product := &model.Product{ID: 1, Name: "Test Product", ProductKey: "P001"}
	mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(product, nil)

	body, _ := json.Marshal(map[string]string{"name": "Test Product"})
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	router := setupProductRouter(mockService)

	product := &model.Product{ID: 1, Name: "Test Product"}
	mockService.EXPECT().Get(gomock.Any(), int64(1)).Return(product, nil)

	req, _ := http.NewRequest("GET", "/products/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_GetByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	router := setupProductRouter(mockService)

	product := &model.Product{ID: 1, Name: "Test Product", ProductKey: "P001"}
	mockService.EXPECT().GetByKey(gomock.Any(), "P001").Return(product, nil)

	req, _ := http.NewRequest("GET", "/products/key/P001", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	router := setupProductRouter(mockService)

	products := []model.Product{{ID: 1, Name: "Product 1"}}
	mockService.EXPECT().List(gomock.Any(), 1, 10, "", "", "").Return(products, int64(1), nil)

	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	router := setupProductRouter(mockService)

	mockService.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/products/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	router := setupProductRouter(mockService)

	product := &model.Product{ID: 1, Name: "Restored Product"}
	mockService.EXPECT().Restore(gomock.Any(), int64(1)).Return(product, nil)

	req, _ := http.NewRequest("POST", "/products/1/restore", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
