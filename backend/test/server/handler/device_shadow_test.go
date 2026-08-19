package handler_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
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

func setupDeviceRouterRemaining(mockService *mock_service.MockDeviceServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	config.Set("data.kafka.device.brokers", []string{"localhost:9092"})
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)

	router.PUT("/devices/:id/shadow/desired", deviceHandler.Desired)
	router.PUT("/devices/:id/shadow/reported", deviceHandler.Reported)
	router.DELETE("/devices/:id/shadow/desired", deviceHandler.ClearDesired)
	router.GET("/devices/:id/shadow/history", deviceHandler.History)
	router.POST("/devices/batch-upload", deviceHandler.BatchUpload)
	router.POST("/devices/mock-kafka", deviceHandler.MockKafka)
	router.GET("/devices/:id/push-records", deviceHandler.PushRecords)
	router.GET("/devices/:id/push-records/:pushRecordId", deviceHandler.PushRecord)

	return router
}

func TestDeviceHandler_Desired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterRemaining(mockService)

	shadow := &model.DeviceShadow{ID: 1, DeviceID: 1}
	mockService.EXPECT().MutateShadow(gomock.Any(), "D001", int64(1), "app", gomock.Any(), gomock.Any(), false).Return(shadow, nil)

	body, _ := json.Marshal(map[string]interface{}{"version": 1, "desired": map[string]any{"temperature": 25}})
	req, _ := http.NewRequest("PUT", "/devices/D001/shadow/desired", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_Reported(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterRemaining(mockService)

	shadow := &model.DeviceShadow{ID: 1, DeviceID: 1}
	mockService.EXPECT().MutateShadow(gomock.Any(), "D001", int64(1), "device", gomock.Any(), gomock.Any(), false).Return(shadow, nil)

	body, _ := json.Marshal(map[string]interface{}{"version": 1, "reported": map[string]any{"status": "online"}})
	req, _ := http.NewRequest("PUT", "/devices/D001/shadow/reported", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_ClearDesired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterRemaining(mockService)

	shadow := &model.DeviceShadow{ID: 1, DeviceID: 1}
	mockService.EXPECT().MutateShadow(gomock.Any(), "D001", int64(1), "app", gomock.Any(), gomock.Any(), true).Return(shadow, nil)

	body, _ := json.Marshal(map[string]interface{}{"version": 1})
	req, _ := http.NewRequest("DELETE", "/devices/D001/shadow/desired", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_History(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterRemaining(mockService)

	history := []model.DeviceShadowHistory{{ID: 1, DeviceID: 1, Version: 1, Source: "device"}}
	mockService.EXPECT().ShadowHistory(gomock.Any(), "D001").Return(history, nil)

	req, _ := http.NewRequest("GET", "/devices/D001/shadow/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_BatchUpload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterRemaining(mockService)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "devices.xlsx")
	part.Write([]byte("test data"))
	writer.Close()

	mockService.EXPECT().BatchCreate(gomock.Any(), gomock.Any()).Return(10, []service.BatchUploadError{}, nil)

	req, _ := http.NewRequest("POST", "/devices/batch-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_BatchUploadNoFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterRemaining(mockService)

	req, _ := http.NewRequest("POST", "/devices/batch-upload", nil)
	req.Header.Set("Content-Type", "multipart/form-data")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceHandler_BatchUploadWithErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterRemaining(mockService)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "devices.xlsx")
	part.Write([]byte("test data"))
	writer.Close()

	errors := []service.BatchUploadError{
		{Row: 2, ProductKey: "P001", DeviceName: "Device2", Error: "duplicate key"},
	}
	mockService.EXPECT().BatchCreate(gomock.Any(), gomock.Any()).Return(9, errors, nil)

	req, _ := http.NewRequest("POST", "/devices/batch-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_MockKafka(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterRemaining(mockService)

	mockService.EXPECT().MockKafka(gomock.Any(), gomock.Any(), "test-topic", `{"key":"value"}`).Return(nil)

	body, _ := json.Marshal(map[string]interface{}{
		"topic": "test-topic",
		"data":  `{"key":"value"}`,
	})
	req, _ := http.NewRequest("POST", "/devices/mock-kafka", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_PushRecordsWithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterRemaining(mockService)

	records := []model.DevicePushRecord{{ID: 1, DeviceID: 1}}
	mockService.EXPECT().ListPushRecords(gomock.Any(), "D001", 1, 20, "", "success").Return(records, int64(1), nil)

	req, _ := http.NewRequest("GET", "/devices/D001/push-records?status=success", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_PushRecordGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterRemaining(mockService)

	record := &model.DevicePushRecord{ID: 1, DeviceID: 1}
	mockService.EXPECT().PushRecord(gomock.Any(), int64(1)).Return(record, nil)

	req, _ := http.NewRequest("GET", "/devices/D001/push-records/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_BatchUploadWithServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterRemaining(mockService)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "devices.xlsx")
	part.Write([]byte("test"))
	writer.Close()

	mockService.EXPECT().BatchCreate(gomock.Any(), gomock.Any()).Return(0, nil, assert.AnError)

	req, _ := http.NewRequest("POST", "/devices/batch-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
