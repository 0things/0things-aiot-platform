package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"0things-backend/internal/handler"
	"0things-backend/internal/model"
	mock_service "0things-backend/test/mocks/service"
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
	router.GET("/products/:id", productHandler.Get)
	router.GET("/products/key/:productKey", productHandler.GetByKey)
	router.GET("/products", productHandler.List)
	router.PUT("/products/key/:productKey", productHandler.Update)
	router.DELETE("/products/:id", productHandler.Delete)
	router.POST("/products/:id/restore", productHandler.Restore)

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
	req, _ := http.NewRequest("PUT", "/products/key/P001", bytes.NewBuffer(body))
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
	req, _ := http.NewRequest("PUT", "/products/key/P001", bytes.NewBuffer(body))
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

	req, _ := http.NewRequest("PUT", "/products/key/P001", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductHandler_GetByKeyNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	router := setupProductRouterFull(mockService)

	mockService.EXPECT().GetByKey(gomock.Any(), "NONEXIST").Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/products/key/NONEXIST", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
