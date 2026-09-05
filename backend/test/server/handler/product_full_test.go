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

func setupProductRouterFull(mockService *mock_service.MockProductServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	productHandler := handler.NewProductHandler(h, mockService)

	router.POST("/products", productHandler.Create)
	router.GET("/products/:productKey", productHandler.Get)
	router.GET("/products", productHandler.List)
	router.PUT("/products/:productKey", productHandler.Update)
	router.DELETE("/products/:productKey", productHandler.Delete)

	return router
}

func TestProductHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	router := setupProductRouterFull(mockService)

	existingProduct := &model.Product{ID: 1, ProductKey: "P001", Name: "Old Name"}
	mockService.EXPECT().GetByKey(gomock.Any(), "P001").Return(existingProduct, nil)
	mockService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]interface{}{
		"name":        "New Name",
		"description": "New Description",
		"category":    "sensor",
		"status":      "active",
	})
	req, _ := http.NewRequest("PUT", "/products/P001", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_UpdatePartial(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	router := setupProductRouterFull(mockService)

	existingProduct := &model.Product{ID: 1, ProductKey: "P001", Name: "Old Name"}
	mockService.EXPECT().GetByKey(gomock.Any(), "P001").Return(existingProduct, nil)
	mockService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "New Name",
	})
	req, _ := http.NewRequest("PUT", "/products/P001", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_UpdateInvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	router := setupProductRouterFull(mockService)

	existingProduct := &model.Product{ID: 1, ProductKey: "P001"}
	mockService.EXPECT().GetByKey(gomock.Any(), "P001").Return(existingProduct, nil)

	req, _ := http.NewRequest("PUT", "/products/P001", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
