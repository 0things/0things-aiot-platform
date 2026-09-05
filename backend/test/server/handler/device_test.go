package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"aiot-backend/internal/handler"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"
	mock_service "aiot-backend/test/mocks/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupDeviceRouter(mockService *mock_service.MockDeviceServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)

	router.POST("/devices", deviceHandler.CreateDevice)
	router.GET("/devices/:deviceKey", deviceHandler.GetDevice)
	router.GET("/devices/key/:deviceKey", deviceHandler.GetDeviceByKey)
	router.GET("/devices", deviceHandler.ListDevices)
	router.PUT("/devices/:deviceKey", deviceHandler.UpdateDevice)
	router.DELETE("/devices/:deviceKey", deviceHandler.DeleteDevice)
	router.POST("/devices/:deviceKey/activate", deviceHandler.Activate)
	router.POST("/devices/:deviceKey/enabled", deviceHandler.Enabled)
	router.GET("/devices/stats", deviceHandler.Stats)

	return router
}

func TestDeviceHandler_CreateDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouter(mockService)

	device := &model.Device{ID: 1, Name: "Test Device", DeviceKey: "D001"}
	mockService.EXPECT().CreateDevice(gomock.Any(), gomock.Any()).Return(device, nil)

	body, _ := json.Marshal(map[string]interface{}{"name": "Test Device", "productId": 1})
	req, _ := http.NewRequest("POST", "/devices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_GetDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouter(mockService)

	device := &model.Device{ID: 1, Name: "Test Device"}
	mockService.EXPECT().DeviceByKey(gomock.Any(), "1").Return(device, nil)

	req, _ := http.NewRequest("GET", "/devices/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_ListDevices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouter(mockService)

	devices := []model.Device{{ID: 1, Name: "Device 1"}}
	mockService.EXPECT().ListDevices(gomock.Any(), gomock.Any()).Return(devices, int64(1), nil)

	req, _ := http.NewRequest("GET", "/devices", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_DeleteDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouter(mockService)

	mockService.EXPECT().DeleteDeviceByKey(gomock.Any(), "1").Return(nil)

	req, _ := http.NewRequest("DELETE", "/devices/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_Stats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouter(mockService)

	stats := service.DeviceStatistics{TotalDevices: 10, OnlineDevices: 5}
	mockService.EXPECT().Stats(gomock.Any()).Return(stats, nil)

	req, _ := http.NewRequest("GET", "/devices/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
