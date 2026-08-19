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

func setupDeviceRouterFull(mockService *mock_service.MockDeviceServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	config := viper.New()
	deviceHandler := handler.NewDeviceHandler(h, mockService, config)

	router.POST("/devices", deviceHandler.CreateDevice)
	router.GET("/devices/:id", deviceHandler.GetDevice)
	router.GET("/devices/key/:deviceKey", deviceHandler.GetDeviceByKey)
	router.GET("/devices", deviceHandler.ListDevices)
	router.PUT("/devices/:id", deviceHandler.UpdateDevice)
	router.DELETE("/devices/:id", deviceHandler.DeleteDevice)
	router.POST("/devices/:id/activate", deviceHandler.Activate)
	router.PUT("/devices/:id/enabled", deviceHandler.Enabled)
	router.GET("/devices/stats", deviceHandler.Stats)
	router.GET("/devices/:id/telemetry", deviceHandler.Telemetry)
	router.GET("/devices/:id/mqtt-parameters", deviceHandler.MQTT)
	router.POST("/devices/:id/restore", deviceHandler.Restore)
	router.GET("/devices/:id/tags", deviceHandler.GetTags)
	router.PUT("/devices/:id/tags", deviceHandler.PutTags)
	router.POST("/devices/:id/tags", deviceHandler.PostTags)
	router.DELETE("/devices/:id/tags", deviceHandler.DeleteTags)
	router.GET("/devices/:id/shadow", deviceHandler.GetShadow)
	router.PUT("/devices/:id/shadow/desired", deviceHandler.Desired)
	router.PUT("/devices/:id/shadow/reported", deviceHandler.Reported)
	router.DELETE("/devices/:id/shadow/desired", deviceHandler.ClearDesired)
	router.GET("/devices/:id/shadow/history", deviceHandler.History)
	router.POST("/devices/:id/simulate-push", deviceHandler.SimulatePush)
	router.GET("/devices/:id/push-records", deviceHandler.PushRecords)
	router.GET("/devices/:id/push-records/:pushRecordId", deviceHandler.PushRecord)
	router.DELETE("/devices/:id/push-records", deviceHandler.ClearPushRecords)
	router.GET("/devices/batch-template", deviceHandler.BatchTemplate)
	router.POST("/devices/batch-upload", deviceHandler.BatchUpload)
	router.POST("/devices/mock-kafka", deviceHandler.MockKafka)

	return router
}

func TestDeviceHandler_GetDeviceByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	device := &model.Device{ID: 1, Name: "Test Device", DeviceKey: "D001"}
	mockService.EXPECT().DeviceByKey(gomock.Any(), "D001").Return(device, nil)

	req, _ := http.NewRequest("GET", "/devices/key/D001", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_UpdateDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	device := &model.Device{ID: 1, Name: "Updated Device"}
	mockService.EXPECT().UpdateDevice(gomock.Any(), int64(1), "Updated Device", "", []byte(nil)).Return(device, nil)

	body, _ := json.Marshal(map[string]string{"name": "Updated Device"})
	req, _ := http.NewRequest("PUT", "/devices/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_Activate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	device := &model.Device{ID: 1, Name: "Activated Device"}
	mockService.EXPECT().Activate(gomock.Any(), int64(1)).Return(device, nil)

	req, _ := http.NewRequest("POST", "/devices/1/activate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_Enabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	device := &model.Device{ID: 1, Name: "Enabled Device", Enabled: true}
	mockService.EXPECT().SetEnabled(gomock.Any(), int64(1), true).Return(device, nil)

	body, _ := json.Marshal(map[string]bool{"enabled": true})
	req, _ := http.NewRequest("PUT", "/devices/1/enabled", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_Telemetry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	mockService.EXPECT().Telemetry(gomock.Any(), "D001").Return(`{"temperature": 25.5}`, nil)

	req, _ := http.NewRequest("GET", "/devices/D001/telemetry", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_MQTT(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	mqttParams := service.MQTTParameters{
		ClientID:    "client1",
		Username:    "user1",
		MQTTHostURL: "mqtt.example.com",
		Password:    "pass1",
		Port:        8883,
	}
	mockService.EXPECT().MQTT(gomock.Any(), "D001").Return(mqttParams, nil)

	req, _ := http.NewRequest("GET", "/devices/D001/mqtt-parameters", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_GetTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	tags := []model.DeviceTag{{ID: 1, Key: "region", Value: "us-east"}}
	mockService.EXPECT().Tags(gomock.Any(), "D001").Return(tags, nil)

	req, _ := http.NewRequest("GET", "/devices/D001/tags", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_PutTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	tags := []model.DeviceTag{{ID: 1, Key: "region", Value: "us-east"}}
	mockService.EXPECT().SetTags(gomock.Any(), "D001", map[string]string{"region": "us-east"}, true).Return(tags, nil)

	body, _ := json.Marshal(map[string]interface{}{"tags": map[string]string{"region": "us-east"}})
	req, _ := http.NewRequest("PUT", "/devices/D001/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_PostTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	tags := []model.DeviceTag{{ID: 1, Key: "region", Value: "us-east"}}
	mockService.EXPECT().SetTags(gomock.Any(), "D001", map[string]string{"region": "us-east"}, false).Return(tags, nil)

	body, _ := json.Marshal(map[string]interface{}{"tags": map[string]string{"region": "us-east"}})
	req, _ := http.NewRequest("POST", "/devices/D001/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_DeleteTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	mockService.EXPECT().RemoveTags(gomock.Any(), "D001", []string{"region"}).Return(nil)

	body, _ := json.Marshal(map[string]interface{}{"keys": []string{"region"}})
	req, _ := http.NewRequest("DELETE", "/devices/D001/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_GetShadow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	shadow := &model.DeviceShadow{ID: 1, DeviceID: 1}
	mockService.EXPECT().Shadow(gomock.Any(), "D001").Return(shadow, nil)

	req, _ := http.NewRequest("GET", "/devices/D001/shadow", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_SimulatePush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	record := &model.DevicePushRecord{ID: 1, DeviceID: 1, Status: "success"}
	mockService.EXPECT().SimulatePush(gomock.Any(), "D001", `{"temp": 25}`, "").Return(record, nil)

	body, _ := json.Marshal(map[string]string{"payload": `{"temp": 25}`})
	req, _ := http.NewRequest("POST", "/devices/D001/simulate-push", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_PushRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	records := []model.DevicePushRecord{{ID: 1, DeviceID: 1}}
	mockService.EXPECT().ListPushRecords(gomock.Any(), "D001", 1, 20, "", "").Return(records, int64(1), nil)

	req, _ := http.NewRequest("GET", "/devices/D001/push-records", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_PushRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	record := &model.DevicePushRecord{ID: 1, DeviceID: 1}
	mockService.EXPECT().PushRecord(gomock.Any(), int64(1)).Return(record, nil)

	req, _ := http.NewRequest("GET", "/devices/D001/push-records/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_ClearPushRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	mockService.EXPECT().ClearPushRecords(gomock.Any(), "D001", gomock.Any()).Return(int64(5), nil)

	req, _ := http.NewRequest("DELETE", "/devices/D001/push-records", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_BatchTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	mockService.EXPECT().BatchTemplate().Return([]byte("template"), nil)

	req, _ := http.NewRequest("GET", "/devices/batch-template", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceHandler_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInterface(ctrl)
	router := setupDeviceRouterFull(mockService)

	device := &model.Device{ID: 1, Name: "Restored Device"}
	mockService.EXPECT().RestoreDevice(gomock.Any(), int64(1)).Return(device, nil)

	req, _ := http.NewRequest("POST", "/devices/1/restore", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
