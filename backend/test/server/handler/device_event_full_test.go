package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/handler"
	mock_service "aiot-backend/test/mocks/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupDeviceEventRouterFull(mockService *mock_service.MockDeviceEventServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	eventHandler := handler.NewDeviceEventHandler(h, mockService)

	router.GET("/device-events", eventHandler.ListDeviceEvents)

	return router
}

func TestDeviceEventHandler_ListWithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	router := setupDeviceEventRouterFull(mockService)

	events := []dto.DeviceEventListItem{{ID: 1, EventType: "temperature"}}
	mockService.EXPECT().List(gomock.Any(), gomock.Any()).Return(events, int64(1), nil)

	req, _ := http.NewRequest("GET", "/device-events?keyword=sensor&deviceKey=D001&eventType=temperature", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceEventHandler_ListWithTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	router := setupDeviceEventRouterFull(mockService)

	events := []dto.DeviceEventListItem{{ID: 1, EventType: "temperature"}}
	mockService.EXPECT().List(gomock.Any(), gomock.Any()).Return(events, int64(1), nil)

	req, _ := http.NewRequest("GET", "/device-events?startAt=2026-08-15%2004:00:00&endAt=2026-08-15%2005:00:00", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceEventHandler_ListWithInvalidStartTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	router := setupDeviceEventRouterFull(mockService)

	req, _ := http.NewRequest("GET", "/device-events?startAt=invalid-time", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeviceEventHandler_ListWithInvalidEndTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	router := setupDeviceEventRouterFull(mockService)

	req, _ := http.NewRequest("GET", "/device-events?endAt=invalid-time", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
