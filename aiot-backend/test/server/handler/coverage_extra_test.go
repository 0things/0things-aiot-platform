package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"0things-backend/internal/handler"
	"0things-backend/internal/model"
	"0things-backend/internal/service"
	mock_service "0things-backend/test/mocks/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupDeviceHandlerExtra(t *testing.T) (*mock_service.MockDeviceServiceInterface, *gin.Engine) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	dh := handler.NewDeviceHandler(h, mockService, config)

	router.POST("/devices", dh.CreateDevice)
	router.GET("/devices", dh.ListDevices)
	router.GET("/devices/stats", dh.Stats)
	router.GET("/devices/:id", dh.GetDevice)
	router.PUT("/devices/:id", dh.UpdateDevice)
	router.DELETE("/devices/:id", dh.DeleteDevice)
	router.POST("/devices/:id/activate", dh.Activate)
	router.PUT("/devices/:id/enabled", dh.Enabled)
	router.POST("/devices/:id/restore", dh.Restore)
	router.GET("/devices/:id/push-records", dh.PushRecords)
	router.GET("/devices/:id/push-records/:pushRecordId", dh.PushRecord)
	router.DELETE("/devices/:id/push-records", dh.ClearPushRecords)
	router.POST("/devices/:id/simulate-push", dh.SimulatePush)
	router.POST("/devices/batch-upload", dh.BatchUpload)

	return mockService, router
}

func TestDeviceHandler_Stats_Error(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().Stats(gomock.Any()).Return(service.DeviceStatistics{}, errors.New("db error"))

	req, _ := http.NewRequest("GET", "/devices/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_Stats_Success2(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	stats := service.DeviceStatistics{TotalDevices: 10, ActivatedDevices: 8, OnlineDevices: 5, OfflineDevices: 3, InactiveDevices: 2}
	mockService.EXPECT().Stats(gomock.Any()).Return(stats, nil)

	req, _ := http.NewRequest("GET", "/devices/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_Activate_InvalidID(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	_ = mockService

	req, _ := http.NewRequest("POST", "/devices/abc/activate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_Activate_Error(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().Activate(gomock.Any(), int64(1)).Return(nil, errors.New("device already activated"))

	req, _ := http.NewRequest("POST", "/devices/1/activate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeviceHandler_Restore_InvalidID(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	_ = mockService

	req, _ := http.NewRequest("POST", "/devices/abc/restore", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_Restore_Error(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().RestoreDevice(gomock.Any(), int64(1)).Return(nil, errors.New("device not found"))

	req, _ := http.NewRequest("POST", "/devices/1/restore", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_Restore_Success2(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	device := &model.Device{ID: 1, Name: "Restored"}
	mockService.EXPECT().RestoreDevice(gomock.Any(), int64(1)).Return(device, nil)

	req, _ := http.NewRequest("POST", "/devices/1/restore", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_Enabled_InvalidID(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	_ = mockService

	body, _ := json.Marshal(map[string]bool{"enabled": true})
	req, _ := http.NewRequest("PUT", "/devices/abc/enabled", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_Enabled_InvalidJSON(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	_ = mockService

	req, _ := http.NewRequest("PUT", "/devices/1/enabled", bytes.NewBuffer([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_Enabled_Error(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().SetEnabled(gomock.Any(), int64(1), true).Return(nil, errors.New("not found"))

	body, _ := json.Marshal(map[string]bool{"enabled": true})
	req, _ := http.NewRequest("PUT", "/devices/1/enabled", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_UpdateDevice_InvalidID(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	_ = mockService

	body, _ := json.Marshal(map[string]string{"name": "test"})
	req, _ := http.NewRequest("PUT", "/devices/abc", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_UpdateDevice_Error(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().UpdateDevice(gomock.Any(), int64(1), "new", "", gomock.Any()).Return(nil, errors.New("invalid status transition"))

	body, _ := json.Marshal(map[string]string{"name": "new"})
	req, _ := http.NewRequest("PUT", "/devices/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeviceHandler_DeleteDevice_InvalidID(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	_ = mockService

	req, _ := http.NewRequest("DELETE", "/devices/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_DeleteDevice_Error(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().DeleteDevice(gomock.Any(), int64(1)).Return(errors.New("not found"))

	req, _ := http.NewRequest("DELETE", "/devices/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_PushRecords_Error(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().ListPushRecords(gomock.Any(), "1", 1, 20, "", "").Return(nil, int64(0), errors.New("not found"))

	req, _ := http.NewRequest("GET", "/devices/1/push-records", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_PushRecord_InvalidPushRecordID(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	_ = mockService

	req, _ := http.NewRequest("GET", "/devices/1/push-records/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_PushRecord_Error(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().PushRecord(gomock.Any(), int64(1)).Return(nil, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/devices/1/push-records/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_ClearPushRecords_Error(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().ClearPushRecords(gomock.Any(), "1", gomock.Any()).Return(int64(0), errors.New("not found"))

	req, _ := http.NewRequest("DELETE", "/devices/1/push-records", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_ClearPushRecords_InvalidTimestamp(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().ClearPushRecords(gomock.Any(), "1", gomock.Any()).Return(int64(5), nil)

	req, _ := http.NewRequest("DELETE", "/devices/1/push-records?beforeTimestamp=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_ClearPushRecords_NegativeTimestamp(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().ClearPushRecords(gomock.Any(), "1", gomock.Any()).Return(int64(5), nil)

	req, _ := http.NewRequest("DELETE", "/devices/1/push-records?beforeTimestamp=-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_ClearPushRecords_Success2(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().ClearPushRecords(gomock.Any(), "1", gomock.Any()).Return(int64(10), nil)

	req, _ := http.NewRequest("DELETE", "/devices/1/push-records?beforeTimestamp=1700000000000", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_SimulatePush_InvalidJSON2(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	_ = mockService

	req, _ := http.NewRequest("POST", "/devices/1/simulate-push", bytes.NewBuffer([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_SimulatePush_Error(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().SimulatePush(gomock.Any(), "NONEXIST", "{}", "").Return(nil, errors.New("device not found"))

	body, _ := json.Marshal(map[string]string{"payload": "{}"})
	req, _ := http.NewRequest("POST", "/devices/NONEXIST/simulate-push", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_BatchUpload_NoFile(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	_ = mockService

	req, _ := http.NewRequest("POST", "/devices/batch-upload", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_ListDevices_Error(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	mockService.EXPECT().ListDevices(gomock.Any(), 1, 10, int64(0), gomock.Any(), gomock.Any(), "").Return(nil, int64(0), errors.New("db error"))

	req, _ := http.NewRequest("GET", "/devices", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_CreateDevice_Success(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	device := &model.Device{ID: 1, Name: "New Device"}
	mockService.EXPECT().CreateDevice(gomock.Any(), gomock.Any()).Return(device, nil)

	body, _ := json.Marshal(map[string]any{"name": "New Device", "productId": 1})
	req, _ := http.NewRequest("POST", "/devices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_GetDevice_Success(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	device := &model.Device{ID: 1, Name: "Device"}
	mockService.EXPECT().Device(gomock.Any(), int64(1)).Return(device, nil)

	req, _ := http.NewRequest("GET", "/devices/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_UpdateDevice_Success(t *testing.T) {
	mockService, router := setupDeviceHandlerExtra(t)
	device := &model.Device{ID: 1, Name: "Updated"}
	mockService.EXPECT().UpdateDevice(gomock.Any(), int64(1), "Updated", "", gomock.Any()).Return(device, nil)

	body, _ := json.Marshal(map[string]string{"name": "Updated"})
	req, _ := http.NewRequest("PUT", "/devices/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
