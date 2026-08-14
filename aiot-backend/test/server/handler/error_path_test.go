package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"0things-backend/internal/handler"
	"0things-backend/internal/model"
	"0things-backend/internal/repository"
	"0things-backend/internal/service"
	mock_service "0things-backend/test/mocks/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDeviceHandler_CreateDevice_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.POST("/devices", deviceHandler.CreateDevice)

	req, _ := http.NewRequest("POST", "/devices", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_GetDevice_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices/:id", deviceHandler.GetDevice)

	mockService.EXPECT().Device(gomock.Any(), int64(1)).Return(nil, repository.ErrNotFound)
	req, _ := http.NewRequest("GET", "/devices/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeviceHandler_GetDevice_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices/:id", deviceHandler.GetDevice)

	req, _ := http.NewRequest("GET", "/devices/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_GetDeviceByKey_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices/key/:deviceKey", deviceHandler.GetDeviceByKey)

	mockService.EXPECT().DeviceByKey(gomock.Any(), "NONEXIST").Return(nil, repository.ErrNotFound)
	req, _ := http.NewRequest("GET", "/devices/key/NONEXIST", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeviceHandler_ListDevices_WithProductID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices", deviceHandler.ListDevices)

	devices := []model.Device{{ID: 1, Name: "Device1"}}
	mockService.EXPECT().ListDevices(gomock.Any(), 1, 10, int64(1), []string(nil), (*bool)(nil), "").Return(devices, int64(1), nil)
	req, _ := http.NewRequest("GET", "/devices?productId=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_ListDevices_WithStates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices", deviceHandler.ListDevices)

	devices := []model.Device{{ID: 1, Name: "Device1"}}
	mockService.EXPECT().ListDevices(gomock.Any(), 1, 10, int64(0), []string{"online", "offline"}, (*bool)(nil), "").Return(devices, int64(1), nil)
	req, _ := http.NewRequest("GET", "/devices?states=online&states=offline", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_ListDevices_WithEnabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices", deviceHandler.ListDevices)

	devices := []model.Device{{ID: 1, Name: "Device1"}}
	enabled := true
	mockService.EXPECT().ListDevices(gomock.Any(), 1, 10, int64(0), []string(nil), &enabled, "").Return(devices, int64(1), nil)
	req, _ := http.NewRequest("GET", "/devices?enabled=true", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_ListDevices_WithSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices", deviceHandler.ListDevices)

	devices := []model.Device{{ID: 1, Name: "Device1"}}
	mockService.EXPECT().ListDevices(gomock.Any(), 1, 10, int64(0), []string(nil), (*bool)(nil), "test").Return(devices, int64(1), nil)
	req, _ := http.NewRequest("GET", "/devices?searchText=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_UpdateDevice_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.PUT("/devices/:id", deviceHandler.UpdateDevice)

	req, _ := http.NewRequest("PUT", "/devices/1", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_DeleteDevice_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.DELETE("/devices/:id", deviceHandler.DeleteDevice)

	mockService.EXPECT().DeleteDevice(gomock.Any(), int64(1)).Return(nil)
	req, _ := http.NewRequest("DELETE", "/devices/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_Stats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices/stats", deviceHandler.Stats)

	stats := service.DeviceStatistics{TotalDevices: 10, OnlineDevices: 5}
	mockService.EXPECT().Stats(gomock.Any()).Return(stats, nil)
	req, _ := http.NewRequest("GET", "/devices/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_Telemetry_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices/:id/telemetry", deviceHandler.Telemetry)

	mockService.EXPECT().Telemetry(gomock.Any(), "D001").Return("", repository.ErrNotFound)
	req, _ := http.NewRequest("GET", "/devices/D001/telemetry", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeviceHandler_MQTT_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices/:id/mqtt", deviceHandler.MQTT)

	mockService.EXPECT().MQTT(gomock.Any(), "D001").Return(service.MQTTParameters{}, repository.ErrNotFound)
	req, _ := http.NewRequest("GET", "/devices/D001/mqtt", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeviceHandler_Restore_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.POST("/devices/:id/restore", deviceHandler.Restore)

	device := &model.Device{ID: 1}
	mockService.EXPECT().RestoreDevice(gomock.Any(), int64(1)).Return(device, nil)
	req, _ := http.NewRequest("POST", "/devices/1/restore", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_GetTags_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices/:id/tags", deviceHandler.GetTags)

	mockService.EXPECT().Tags(gomock.Any(), "D001").Return(nil, repository.ErrNotFound)
	req, _ := http.NewRequest("GET", "/devices/D001/tags", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeviceHandler_PutTags_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.PUT("/devices/:id/tags", deviceHandler.PutTags)

	req, _ := http.NewRequest("PUT", "/devices/D001/tags", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_DeleteTags_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.DELETE("/devices/:id/tags", deviceHandler.DeleteTags)

	mockService.EXPECT().RemoveTags(gomock.Any(), "D001", []string{"key1"}).Return(nil)
	body, _ := json.Marshal(map[string]interface{}{"keys": []string{"key1"}})
	req, _ := http.NewRequest("DELETE", "/devices/D001/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_DeleteTags_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.DELETE("/devices/:id/tags", deviceHandler.DeleteTags)

	req, _ := http.NewRequest("DELETE", "/devices/D001/tags", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_GetShadow_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices/:id/shadow", deviceHandler.GetShadow)

	mockService.EXPECT().Shadow(gomock.Any(), "D001").Return(nil, repository.ErrNotFound)
	req, _ := http.NewRequest("GET", "/devices/D001/shadow", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeviceHandler_History_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.GET("/devices/:id/shadow/history", deviceHandler.History)

	history := []model.DeviceShadowHistory{{ID: 1, DeviceID: 1, Version: 1, Source: "device"}}
	mockService.EXPECT().ShadowHistory(gomock.Any(), "D001").Return(history, nil)
	req, _ := http.NewRequest("GET", "/devices/D001/shadow/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_ClearPushRecords_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.DELETE("/devices/:id/push-records", deviceHandler.ClearPushRecords)

	mockService.EXPECT().ClearPushRecords(gomock.Any(), "D001", gomock.Any()).Return(int64(5), nil)
	req, _ := http.NewRequest("DELETE", "/devices/D001/push-records", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_ClearPushRecords_WithBeforeTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.DELETE("/devices/:id/push-records", deviceHandler.ClearPushRecords)

	before := time.Now().Format(time.RFC3339)
	mockService.EXPECT().ClearPushRecords(gomock.Any(), "D001", gomock.Any()).Return(int64(5), nil)
	req, _ := http.NewRequest("DELETE", "/devices/D001/push-records?before="+before, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_MockKafka_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.POST("/devices/mock-kafka", deviceHandler.MockKafka)

	req, _ := http.NewRequest("POST", "/devices/mock-kafka", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_MockKafka_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	config.Set("data.kafka.device.brokers", []string{"localhost:9092"})
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.POST("/devices/mock-kafka", deviceHandler.MockKafka)

	mockService.EXPECT().MockKafka(gomock.Any(), gomock.Any(), "test-topic", `{"key":"value"}`).Return(nil)
	body, _ := json.Marshal(map[string]interface{}{"topic": "test-topic", "data": `{"key":"value"}`})
	req, _ := http.NewRequest("POST", "/devices/mock-kafka", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_SimulatePush_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)
	router.POST("/devices/:id/simulate-push", deviceHandler.SimulatePush)

	req, _ := http.NewRequest("POST", "/devices/D001/simulate-push", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
