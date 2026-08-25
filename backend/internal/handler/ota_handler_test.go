package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"aiot-backend/internal/model"
	"aiot-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_service "aiot-backend/test/mocks/service"
	"github.com/stretchr/testify/require"
	"errors"
)

func newOTATestRouter(h *OTAHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ota-packages", h.ListOTA)
	r.GET("/ota-packages/:uuid", h.GetOTA)
	r.POST("/ota-packages", h.CreateOTA)
	r.PUT("/ota-packages/:uuid", h.UpdateOTA)
	r.DELETE("/ota-packages/:uuid", h.DeleteOTA)
	r.POST("/ota-packages/:uuid/deploy", h.DeployOTA)
	r.POST("/ota-packages/:uuid/dispatch", h.DispatchOTA)
	r.POST("/ota-packages/:uuid/report", h.ReportOTAStatus)
	r.GET("/ota-packages/:uuid/upgrade-statistics", h.OTAStats)
	r.GET("/ota-packages/:uuid/batches", h.OTABatches)
	r.GET("/ota-packages/:uuid/device-deployments", h.OTADeployments)
	return r
}

func TestOTAHandler_ListOTA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.OTAPackage{{ID: 1}}, int64(1), nil)

	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodGet, "/ota-packages", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"total":1`)
}

func TestOTAHandler_GetOTA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().Get(gomock.Any(), "1").Return(&model.OTAPackage{ID: 1}, nil)

	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodGet, "/ota-packages/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"id":1`)
}

func TestOTAHandler_GetOTA_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().Get(gomock.Any(), "abc").Return(nil, errors.New("not found"))

	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodGet, "/ota-packages/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOTAHandler_CreateOTA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]any{"packageName": "p1", "version": "1.0", "product_key": "pk"})
	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodPost, "/ota-packages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_CreateOTA_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)

	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodPost, "/ota-packages", bytes.NewReader([]byte("{bad")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOTAHandler_UpdateOTA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().Get(gomock.Any(), "1").Return(&model.OTAPackage{ID: 1}, nil)
	mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]any{"packageName": "p2"})
	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodPut, "/ota-packages/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_DeleteOTA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().Delete(gomock.Any(), "1").Return(nil)

	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodDelete, "/ota-packages/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_DeployOTA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().Deploy(gomock.Any(), "1", gomock.Any()).Return(1, nil)

	body, _ := json.Marshal(map[string]any{"deviceKeys": []string{"d1"}})
	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodPost, "/ota-packages/1/deploy", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_DeployOTA_EmptyKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)

	body, _ := json.Marshal(map[string]any{"deviceKeys": []string{}})
	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodPost, "/ota-packages/1/deploy", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOTAHandler_DispatchOTA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().Dispatch(gomock.Any(), "1").Return(int64(1), nil)

	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodPost, "/ota-packages/1/dispatch", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_ReportOTAStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().ReportStatus(gomock.Any(), "1", "d1", "success").Return(nil)

	body, _ := json.Marshal(map[string]any{"deviceKey": "d1", "status": "success"})
	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodPost, "/ota-packages/1/report", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_OTAStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().Statistics(gomock.Any(), "1").Return(service.UpgradeStatistics{PackageID: "1", TotalTargetDevices: 2}, nil)

	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodGet, "/ota-packages/1/upgrade-statistics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"totalTargetDevices":2`)
}

func TestOTAHandler_OTABatches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().Batches(gomock.Any(), "1").Return([]model.UpgradeBatch{{BatchID: "b1"}}, nil)

	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodGet, "/ota-packages/1/batches", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_OTADeployments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_service.NewMockOTAServiceInterface(ctrl)
	mock.EXPECT().Deployments(gomock.Any(), "1", gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.DeviceDeployment{{DeviceID: 1}}, int64(1), nil)

	r := newOTATestRouter(NewOTAHandler(&Handler{}, mock))
	req := httptest.NewRequest(http.MethodGet, "/ota-packages/1/device-deployments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
