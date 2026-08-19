package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"aiot-backend/internal/handler"
	"aiot-backend/internal/model"
	mock_service "aiot-backend/test/mocks/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupAlertRouter(mockService *mock_service.MockAlertServiceInterface) *gin.Engine {
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

func TestAlertHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockAlertServiceInterface(ctrl)
	router := setupAlertRouter(mockService)

	alerts := []model.Alert{{ID: 1, Summary: "High Temperature"}}
	mockService.EXPECT().List(gomock.Any(), 1, 20, "", "", "").Return(alerts, int64(1), nil)

	req, _ := http.NewRequest("GET", "/alerts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAlertHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockAlertServiceInterface(ctrl)
	router := setupAlertRouter(mockService)

	alert := &model.Alert{ID: 1, Summary: "High Temperature"}
	mockService.EXPECT().Get(gomock.Any(), int64(1)).Return(alert, nil)

	req, _ := http.NewRequest("GET", "/alerts/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAlertHandler_Ack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockAlertServiceInterface(ctrl)
	router := setupAlertRouter(mockService)

	alert := &model.Alert{ID: 1, Summary: "High Temperature", Status: "acknowledged"}
	mockService.EXPECT().SetStatus(gomock.Any(), int64(1), "acknowledged").Return(alert, nil)

	req, _ := http.NewRequest("POST", "/alerts/1/ack", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAlertHandler_Resolve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockAlertServiceInterface(ctrl)
	router := setupAlertRouter(mockService)

	alert := &model.Alert{ID: 1, Summary: "High Temperature", Status: "resolved"}
	mockService.EXPECT().SetStatus(gomock.Any(), int64(1), "resolved").Return(alert, nil)

	req, _ := http.NewRequest("POST", "/alerts/1/resolve", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
