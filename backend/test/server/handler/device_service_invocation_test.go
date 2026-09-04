package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"aiot-backend/internal/handler"
	"aiot-backend/internal/model"
	mock_service "aiot-backend/test/mocks/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupDeviceServiceInvocationRouter(mockService *mock_service.MockDeviceServiceInvocationServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	invocationHandler := handler.NewDeviceServiceInvocationHandler(h, mockService)

	router.GET("/devices/:deviceKey/thing-model-service-invocations", invocationHandler.List)

	return router
}

func TestDeviceServiceInvocationHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInvocationServiceInterface(ctrl)
	router := setupDeviceServiceInvocationRouter(mockService)

	output := `{"status":"ok"}`
	invocations := []model.DeviceServiceInvocation{
		{
			UUID:              "test-uuid-1",
			DeviceID:          1,
			ServiceIdentifier: "reboot",
			ServiceName:       "Reboot Device",
			InputParams:       `{"force":true}`,
			OutputParams:      &output,
			InvokedAt:         time.Now(),
		},
	}
	mockService.EXPECT().List(gomock.Any(), gomock.Any()).Return(invocations, int64(1), nil)

	req, _ := http.NewRequest("GET", "/devices/D001/thing-model-service-invocations?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceServiceInvocationHandler_List_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInvocationServiceInterface(ctrl)
	router := setupDeviceServiceInvocationRouter(mockService)

	req, _ := http.NewRequest("GET", "/devices/D001/thing-model-service-invocations?startAt=invalid-time-format", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeviceServiceInvocationHandler_List_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockDeviceServiceInvocationServiceInterface(ctrl)
	router := setupDeviceServiceInvocationRouter(mockService)

	mockService.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, int64(0), errors.New("device not found"))

	req, _ := http.NewRequest("GET", "/devices/NON_EXISTENT/thing-model-service-invocations", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
