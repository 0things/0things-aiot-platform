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
	"github.com/stretchr/testify/assert"
)

func setupOTARouterFull(mockService *mock_service.MockOTAServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := &handler.Handler{}
	otaHandler := handler.NewOTAHandler(h, mockService)

	router.GET("/ota/packages", otaHandler.ListOTA)
	router.GET("/ota/packages/:id", otaHandler.GetOTA)
	router.POST("/ota/packages", otaHandler.CreateOTA)
	router.PUT("/ota/packages/:id", otaHandler.UpdateOTA)
	router.DELETE("/ota/packages/:id", otaHandler.DeleteOTA)
	router.GET("/ota/packages/:id/stats", otaHandler.OTAStats)
	router.GET("/ota/packages/:id/batches", otaHandler.OTABatches)
	router.GET("/ota/packages/:id/deployments", otaHandler.OTADeployments)

	return router
}

func TestOTAHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTARouterFull(mockService)

	existingPkg := &model.OTAPackage{ID: 1, PackageName: "firmware-1"}
	mockService.EXPECT().Get(gomock.Any(), int64(1)).Return(existingPkg, nil)
	mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	body, _ := json.Marshal(map[string]interface{}{
		"packageName": "firmware-1",
		"version":     "2.0.0",
	})
	req, _ := http.NewRequest("PUT", "/ota/packages/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_Stats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTARouterFull(mockService)

	mockService.EXPECT().Statistics(gomock.Any(), "1").Return(
		service.UpgradeStatistics{
			PackageID:          "1",
			TotalTargetDevices: 100,
			SuccessfulUpgrades: 80,
		}, nil)

	req, _ := http.NewRequest("GET", "/ota/packages/1/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_Batches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTARouterFull(mockService)

	batches := []model.UpgradeBatch{{BatchID: "B001", BatchName: "Batch 1"}}
	mockService.EXPECT().Batches(gomock.Any(), "1").Return(batches, nil)

	req, _ := http.NewRequest("GET", "/ota/packages/1/batches", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_Deployments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTARouterFull(mockService)

	deployments := []model.DeviceDeployment{{DeviceKey: "D001"}}
	mockService.EXPECT().Deployments(gomock.Any(), "1", 1, 100, "").Return(deployments, int64(1), nil)

	req, _ := http.NewRequest("GET", "/ota/packages/1/deployments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
