package handler_test

import (
	"bytes"
	"encoding/json"
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

func setupOTARouter(mockService *mock_service.MockOTAServiceInterface) *gin.Engine {
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
	router.POST("/ota/packages/:id/batch-upgrade", otaHandler.BatchUpgradeOTA)

	return router
}

func TestOTAHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTARouter(mockService)

	packages := []model.OTAPackage{{ID: 1, PackageName: "firmware-1"}}
	mockService.EXPECT().List(gomock.Any(), 1, 20).Return(packages, int64(1), nil)

	req, _ := http.NewRequest("GET", "/ota/packages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTARouter(mockService)

	pkg := &model.OTAPackage{ID: 1, PackageName: "firmware-1"}
	mockService.EXPECT().Get(gomock.Any(), int64(1)).Return(pkg, nil)

	req, _ := http.NewRequest("GET", "/ota/packages/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTARouter(mockService)

	mockService.EXPECT().Create(gomock.Any(), gomock.Any(), "P001").Return(nil)

	body, _ := json.Marshal(map[string]string{"packageName": "firmware-1", "product_key": "P001"})
	req, _ := http.NewRequest("POST", "/ota/packages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTARouter(mockService)

	mockService.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/ota/packages/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOTAHandler_BatchUpgrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockOTAServiceInterface(ctrl)
	router := setupOTARouter(mockService)

	mockService.EXPECT().
		BatchUpgrade(gomock.Any(), "firmware-1", []string{"D001", "D002"}).
		Return(&model.UpgradeBatch{BatchID: "B-1-001", Status: "pending"}, nil)

	body, _ := json.Marshal(map[string]any{
		"deviceKeys": []string{"D001", "D002"},
	})
	req, _ := http.NewRequest("POST", "/ota/packages/firmware-1/batch-upgrade", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "B-1-001")
}
