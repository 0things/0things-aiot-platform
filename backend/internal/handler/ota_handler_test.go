package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"aiot-backend/internal/model"
	"aiot-backend/internal/service"
	mock_service "aiot-backend/test/mocks/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func newOTATestRouter(h *OTAHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ota-packages", h.ListOTA)
	r.GET("/ota-packages/:uuid", h.GetOTA)
	r.POST("/ota-packages", h.CreateOTA)
	r.PUT("/ota-packages/:uuid", h.UpdateOTA)
	r.DELETE("/ota-packages/:uuid", h.DeleteOTA)
	r.POST("/ota-packages/:uuid/report", h.ReportOTAStatus)
	r.GET("/ota-packages/:uuid/upgrade-statistics", h.OTAStats)
	r.GET("/ota-packages/:uuid/batches", h.OTABatches)
	r.GET("/ota-packages/:uuid/device-deployments", h.OTADeployments)
	return r
}

func TestOTAPackageJSONFormatsTimesWithCarbon(t *testing.T) {
	createdAt := time.Date(2026, time.August, 25, 9, 30, 0, 0, time.UTC)
	releasedAt := createdAt.Add(time.Hour)

	result := otaPackageJSON(model.OTAPackage{
		CreatedAt:  createdAt,
		UpdatedAt:  createdAt,
		ReleasedAt: &releasedAt,
	})

	require.Equal(t, "2026-08-25 09:30:00", result.CreatedAt)
	require.Equal(t, "2026-08-25 09:30:00", result.UpdatedAt)
	require.NotNil(t, result.ReleasedAt)
	require.Equal(t, "2026-08-25 10:30:00", *result.ReleasedAt)
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
