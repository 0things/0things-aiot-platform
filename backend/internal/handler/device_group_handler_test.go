package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	devicegroupv1 "aiot-backend/api/device_group/v1"
	"aiot-backend/internal/model"
	mock_service "aiot-backend/test/mocks/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func newDeviceGroupTestRouter(h *DeviceGroupHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/device-groups", h.Create)
	r.GET("/device-groups", h.List)
	r.GET("/device-groups/:groupUuid", h.Get)
	r.PUT("/device-groups/:groupUuid", h.Update)
	r.DELETE("/device-groups/:groupUuid", h.Delete)
	r.POST("/device-groups/:groupUuid/devices", h.AddDevices)
	r.DELETE("/device-groups/:groupUuid/devices", h.RemoveDevices)
	r.GET("/device-groups/:groupUuid/devices", h.ListDevices)
	r.POST("/device-groups/preview", h.Preview)
	r.POST("/device-groups/:groupUuid/preview", h.PreviewSaved)
	return r
}

func TestDeviceGroupHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockDeviceGroupServiceInterface(ctrl)
	h := NewDeviceGroupHandler(&Handler{}, mockSvc)
	r := newDeviceGroupTestRouter(h)

	group := &model.DeviceGroup{
		GroupUUID: "uuid-1",
		Name:      "test-group",
		Type:      "manual",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(group, nil)

	body, _ := json.Marshal(devicegroupv1.CreateDeviceGroupRequest{
		Name: "test-group",
		Type: "manual",
	})
	req := httptest.NewRequest(http.MethodPost, "/device-groups", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceGroupHandler_Create_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockDeviceGroupServiceInterface(ctrl)
	h := NewDeviceGroupHandler(&Handler{}, mockSvc)
	r := newDeviceGroupTestRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/device-groups", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeviceGroupHandler_Get_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockDeviceGroupServiceInterface(ctrl)
	mockSvc.EXPECT().Get(gomock.Any(), "err-uuid").Return(nil, errors.New("database failure"))

	h := NewDeviceGroupHandler(&Handler{}, mockSvc)
	r := newDeviceGroupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/device-groups/err-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeviceGroupHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockDeviceGroupServiceInterface(ctrl)
	mockSvc.EXPECT().Delete(gomock.Any(), "uuid-1").Return(nil)

	h := NewDeviceGroupHandler(&Handler{}, mockSvc)
	r := newDeviceGroupTestRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/device-groups/uuid-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceGroupHandler_AddDevices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockDeviceGroupServiceInterface(ctrl)
	mockSvc.EXPECT().AddDevices(gomock.Any(), "uuid-1", []string{"dev-1"}).Return(nil)

	h := NewDeviceGroupHandler(&Handler{}, mockSvc)
	r := newDeviceGroupTestRouter(h)

	body, _ := json.Marshal(devicegroupv1.DeviceKeysRequest{
		DeviceKeys: []string{"dev-1"},
	})
	req := httptest.NewRequest(http.MethodPost, "/device-groups/uuid-1/devices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestDeviceGroupHandler_PreviewSaved_RejectsManual(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockDeviceGroupServiceInterface(ctrl)
	mockSvc.EXPECT().Get(gomock.Any(), "uuid-manual").Return(&model.DeviceGroup{
		GroupUUID: "uuid-manual",
		Type:      "manual",
	}, nil)

	h := NewDeviceGroupHandler(&Handler{}, mockSvc)
	r := newDeviceGroupTestRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/device-groups/uuid-manual/preview", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
