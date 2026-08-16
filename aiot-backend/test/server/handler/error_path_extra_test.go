package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"0things-backend/internal/handler"
	"0things-backend/internal/model"
	"0things-backend/internal/repository"
	"0things-backend/internal/service"
	mock_service "0things-backend/test/mocks/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupOTAFullErrRouter(mockService *mock_service.MockOTAServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	otaHandler := handler.NewOTAHandler(h, mockService)

	router.GET("/ota/packages", otaHandler.ListOTA)
	router.GET("/ota/packages/:id", otaHandler.GetOTA)
	router.POST("/ota/packages", otaHandler.CreateOTA)
	router.PUT("/ota/packages/:id", otaHandler.UpdateOTA)
	router.DELETE("/ota/packages/:id", otaHandler.DeleteOTA)
	router.GET("/ota/packages/:id/stats", otaHandler.OTAStats)
	router.GET("/ota/packages/:id/batches", otaHandler.OTABatches)
	router.GET("/ota/packages/:id/deployments", otaHandler.OTADeployments)

	return router
}

func TestOTAHandler_GetOTA_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	mockService.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, repository.ErrNotFound)
	req, _ := http.NewRequest("GET", "/ota/packages/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOTAHandler_GetOTA_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	req, _ := http.NewRequest("GET", "/ota/packages/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOTAHandler_CreateOTA_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	req, _ := http.NewRequest("POST", "/ota/packages", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOTAHandler_DeleteOTA_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	mockService.EXPECT().Delete(gomock.Any(), int64(1)).Return(repository.ErrNotFound)
	req, _ := http.NewRequest("DELETE", "/ota/packages/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOTAHandler_OTAStats_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	mockService.EXPECT().Statistics(gomock.Any(), "1").Return(service.UpgradeStatistics{}, repository.ErrNotFound)
	req, _ := http.NewRequest("GET", "/ota/packages/1/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOTAHandler_OTABatches_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	mockService.EXPECT().Batches(gomock.Any(), "1").Return([]model.UpgradeBatch{}, nil)
	req, _ := http.NewRequest("GET", "/ota/packages/1/batches", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_OTADeployments_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	mockService.EXPECT().Deployments(gomock.Any(), "1", 1, 100, "").Return([]model.DeviceDeployment{}, int64(0), nil)
	req, _ := http.NewRequest("GET", "/ota/packages/1/deployments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_UpdateOTA_GetFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	mockService.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, repository.ErrNotFound)
	body, _ := json.Marshal(map[string]interface{}{"packageName": "fw", "version": "2.0.0", "productKey": "P001"})
	req, _ := http.NewRequest("PUT", "/ota/packages/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOTAHandler_UpdateOTA_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	req, _ := http.NewRequest("PUT", "/ota/packages/1", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOTAHandler_OTAStats_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	mockService.EXPECT().Statistics(gomock.Any(), "abc").Return(service.UpgradeStatistics{}, assert.AnError)
	req, _ := http.NewRequest("GET", "/ota/packages/abc/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOTAHandler_OTABatches_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	mockService.EXPECT().Batches(gomock.Any(), "abc").Return(nil, assert.AnError)
	req, _ := http.NewRequest("GET", "/ota/packages/abc/batches", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOTAHandler_OTADeployments_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTAFullErrRouter(mockService)

	mockService.EXPECT().Deployments(gomock.Any(), "abc", 1, 100, "").Return(nil, int64(0), assert.AnError)
	req, _ := http.NewRequest("GET", "/ota/packages/abc/deployments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func setupAlertFullErrRouter(mockService *mock_service.MockAlertServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	alertHandler := handler.NewAlertHandler(h, mockService)

	router.GET("/alerts", alertHandler.ListAlerts)
	router.GET("/alerts/:id", alertHandler.GetAlert)
	router.POST("/alerts/:id/ack", alertHandler.AckAlert)
	router.POST("/alerts/:id/resolve", alertHandler.ResolveAlert)

	return router
}

func TestAlertHandler_GetAlert_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockAlertServiceInterface(ctrl)
	router := setupAlertFullErrRouter(mockService)

	mockService.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, repository.ErrNotFound)
	req, _ := http.NewRequest("GET", "/alerts/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAlertHandler_GetAlert_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockAlertServiceInterface(ctrl)
	router := setupAlertFullErrRouter(mockService)

	req, _ := http.NewRequest("GET", "/alerts/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAlertHandler_AckAlert_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockAlertServiceInterface(ctrl)
	router := setupAlertFullErrRouter(mockService)

	mockService.EXPECT().SetStatus(gomock.Any(), int64(1), "acknowledged").Return(nil, repository.ErrNotFound)
	req, _ := http.NewRequest("POST", "/alerts/1/ack", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAlertHandler_ResolveAlert_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockAlertServiceInterface(ctrl)
	router := setupAlertFullErrRouter(mockService)

	mockService.EXPECT().SetStatus(gomock.Any(), int64(1), "resolved").Return(nil, repository.ErrNotFound)
	req, _ := http.NewRequest("POST", "/alerts/1/resolve", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAlertHandler_ListAlerts_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockAlertServiceInterface(ctrl)
	router := setupAlertFullErrRouter(mockService)

	mockService.EXPECT().List(gomock.Any(), 1, 20, "", "", "").Return([]model.Alert{}, int64(0), nil)
	req, _ := http.NewRequest("GET", "/alerts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func setupDeviceEventFullErrRouter(mockService *mock_service.MockDeviceEventServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	eventHandler := handler.NewDeviceEventHandler(h, mockService)
	router.GET("/device-events", eventHandler.ListDeviceEvents)
	return router
}

func TestDeviceEventHandler_ListDeviceEvents_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	router := setupDeviceEventFullErrRouter(mockService)

	mockService.EXPECT().List(gomock.Any(), 1, 20, "", "", "", nil, nil).Return([]model.DeviceEvent{}, int64(0), nil)
	req, _ := http.NewRequest("GET", "/device-events", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceEventHandler_ListDeviceEvents_WithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	router := setupDeviceEventFullErrRouter(mockService)

	mockService.EXPECT().List(gomock.Any(), 1, 20, "test", "D001", "temperature", nil, nil).Return([]model.DeviceEvent{}, int64(0), nil)
	req, _ := http.NewRequest("GET", "/device-events?keyword=test&device_key=D001&event_type=temperature", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceEventHandler_ListDeviceEvents_InvalidStartAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	router := setupDeviceEventFullErrRouter(mockService)

	req, _ := http.NewRequest("GET", "/device-events?start_at=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeviceEventHandler_ListDeviceEvents_InvalidEndAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	router := setupDeviceEventFullErrRouter(mockService)

	req, _ := http.NewRequest("GET", "/device-events?end_at=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
