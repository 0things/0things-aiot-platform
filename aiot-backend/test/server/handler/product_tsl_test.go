package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"0things-backend/internal/handler"
	"0things-backend/internal/model"
	"0things-backend/internal/repository"
	mock_service "0things-backend/test/mocks/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupProductTSLRouter(mockService *mock_service.MockProductTSLServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	productTSLHandler := handler.NewProductTSLHandler(h, mockService)

	router.GET("/products/:id/tsl", productTSLHandler.Get)
	router.PUT("/products/:id/tsl", productTSLHandler.Put)
	router.DELETE("/products/:id/tsl", productTSLHandler.Delete)

	return router
}

func TestProductTSLHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductTSLServiceInterface(ctrl)
	router := setupProductTSLRouter(mockService)

	pid := int64(1)
	tsl := &model.ProductTSL{ID: 1, ProductID: &pid, TSL: `{"properties":[]}`, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockService.EXPECT().Get(gomock.Any(), "P001").Return(tsl, nil)

	req, _ := http.NewRequest("GET", "/products/P001/tsl", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductTSLHandler_Get_NilProductID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductTSLServiceInterface(ctrl)
	router := setupProductTSLRouter(mockService)

	tsl := &model.ProductTSL{ID: 1, ProductID: nil, TSL: `{"properties":[]}`, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockService.EXPECT().Get(gomock.Any(), "P001").Return(tsl, nil)

	req, _ := http.NewRequest("GET", "/products/P001/tsl", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductTSLHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductTSLServiceInterface(ctrl)
	router := setupProductTSLRouter(mockService)

	mockService.EXPECT().Get(gomock.Any(), "NONEXIST").Return(nil, repository.ErrNotFound)

	req, _ := http.NewRequest("GET", "/products/NONEXIST/tsl", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductTSLHandler_Get_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductTSLServiceInterface(ctrl)
	router := setupProductTSLRouter(mockService)

	mockService.EXPECT().Get(gomock.Any(), "P001").Return(nil, errors.New("db error"))

	req, _ := http.NewRequest("GET", "/products/P001/tsl", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductTSLHandler_Put_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductTSLServiceInterface(ctrl)
	router := setupProductTSLRouter(mockService)

	body, _ := json.Marshal(map[string]string{"tsl": `{"properties":[]}`})
	mockService.EXPECT().Upsert(gomock.Any(), "P001", `{"properties":[]}`).Return(nil)

	req, _ := http.NewRequest("PUT", "/products/P001/tsl", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductTSLHandler_Put_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductTSLServiceInterface(ctrl)
	router := setupProductTSLRouter(mockService)

	req, _ := http.NewRequest("PUT", "/products/P001/tsl", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductTSLHandler_Put_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductTSLServiceInterface(ctrl)
	router := setupProductTSLRouter(mockService)

	body, _ := json.Marshal(map[string]string{"tsl": `{"properties":[]}`})
	mockService.EXPECT().Upsert(gomock.Any(), "NONEXIST", `{"properties":[]}`).Return(repository.ErrNotFound)

	req, _ := http.NewRequest("PUT", "/products/NONEXIST/tsl", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductTSLHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductTSLServiceInterface(ctrl)
	router := setupProductTSLRouter(mockService)

	mockService.EXPECT().Delete(gomock.Any(), "P001").Return(nil)

	req, _ := http.NewRequest("DELETE", "/products/P001/tsl", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductTSLHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductTSLServiceInterface(ctrl)
	router := setupProductTSLRouter(mockService)

	mockService.EXPECT().Delete(gomock.Any(), "NONEXIST").Return(repository.ErrNotFound)

	req, _ := http.NewRequest("DELETE", "/products/NONEXIST/tsl", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductTSLHandler_Delete_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockProductTSLServiceInterface(ctrl)
	router := setupProductTSLRouter(mockService)

	mockService.EXPECT().Delete(gomock.Any(), "P001").Return(errors.New("db error"))

	req, _ := http.NewRequest("DELETE", "/products/P001/tsl", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
