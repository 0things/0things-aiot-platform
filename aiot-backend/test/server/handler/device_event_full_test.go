package handler_test

import (
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

	events := []model.DeviceEvent{{ID: 1, EventType: "temperature"}}
	mockService.EXPECT().List(gomock.Any(), 1, 20, "sensor", "D001", "temperature", nil, nil).Return(events, int64(1), nil)

	req, _ := http.NewRequest("GET", "/device-events?keyword=sensor&device_key=D001&event_type=temperature", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceEventHandler_ListWithTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	router := setupDeviceEventRouterFull(mockService)

	events := []model.DeviceEvent{{ID: 1, EventType: "temperature"}}
	mockService.EXPECT().List(gomock.Any(), 1, 20, "", "", "", gomock.Any(), gomock.Any()).Return(events, int64(1), nil)

	req, _ := http.NewRequest("GET", "/device-events?start_at=2026-08-15T04:00:00Z&end_at=2026-08-15T05:00:00Z", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceEventHandler_ListWithInvalidStartTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	router := setupDeviceEventRouterFull(mockService)

	req, _ := http.NewRequest("GET", "/device-events?start_at=invalid-time", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeviceEventHandler_ListWithInvalidEndTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceEventServiceInterface(ctrl)
	router := setupDeviceEventRouterFull(mockService)

	req, _ := http.NewRequest("GET", "/device-events?end_at=invalid-time", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
