package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"aiot-backend/internal/handler"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/service"
	"aiot-backend/pkg/config"
	"aiot-backend/pkg/log"
	mock_service "aiot-backend/test/mocks/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"github.com/stretchr/testify/assert"
)

var coverageLogger *log.Logger

func init() {
	if coverageLogger == nil {
		os.Setenv("APP_CONF", "../../../config/local.yml")
		coverageConf = config.NewConfig("config/local.yml")
		coverageLogger = log.NewLog(coverageConf)
	}
}

var coverageConf *viper.Viper

type coverageRouters struct {
	*testing.T
	router     *gin.Engine
	device     *mock_service.MockDeviceServiceInterface
	product    *mock_service.MockProductServiceInterface
	ota        *mock_service.MockOTAServiceInterface
	deviceEvent *mock_service.MockDeviceEventServiceInterface
}

func newCoverageRouters(t *testing.T) *coverageRouters {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewHandler(coverageLogger)

	cr := &coverageRouters{
		T:           t,
		router:      r,
		device:      mock_service.NewMockDeviceServiceInterface(ctrl),
		product:     mock_service.NewMockProductServiceInterface(ctrl),
		ota:         mock_service.NewMockOTAServiceInterface(ctrl),
		deviceEvent: mock_service.NewMockDeviceEventServiceInterface(ctrl),
	}

	dh := handler.NewDeviceHandler(h, cr.device, coverageConf)
	ph := handler.NewProductHandler(h, cr.product)
	oh := handler.NewOTAHandler(h, cr.ota)
	eh := handler.NewDeviceEventHandler(h, cr.deviceEvent)

	r.GET("/devices/stats", dh.Stats)
	r.GET("/devices", dh.ListDevices)
	r.POST("/devices", dh.CreateDevice)
	r.GET("/devices/:id", dh.GetDevice)
	r.PUT("/devices/:id", dh.UpdateDevice)
	r.DELETE("/devices/:id", dh.DeleteDevice)
	r.POST("/devices/:id/activate", dh.Activate)
	r.PUT("/devices/:id/enabled", dh.Enabled)
	r.POST("/devices/:id/restore", dh.Restore)
	r.GET("/devices/:id/tags", dh.GetTags)
	r.PUT("/devices/:id/tags", dh.PutTags)
	r.POST("/devices/:id/tags", dh.PostTags)
	r.DELETE("/devices/:id/tags", dh.DeleteTags)
	r.GET("/devices/:id/shadow", dh.GetShadow)
	r.PUT("/devices/:id/shadow/desired", dh.Desired)
	r.PUT("/devices/:id/shadow/reported", dh.Reported)
	r.DELETE("/devices/:id/shadow/desired", dh.ClearDesired)
	r.GET("/devices/:id/shadow/history", dh.History)
	r.POST("/devices/:id/simulate-push", dh.SimulatePush)
	r.GET("/devices/:id/push-records", dh.PushRecords)
	r.GET("/devices/:id/push-records/:pushRecordId", dh.PushRecord)
	r.DELETE("/devices/:id/push-records", dh.ClearPushRecords)
	r.GET("/devices/batch-template", dh.BatchTemplate)
	r.POST("/devices/mock-kafka", dh.MockKafka)

	r.POST("/products", ph.Create)
	r.GET("/products/:id", ph.Get)
	r.GET("/products", ph.List)
	r.PUT("/products/key/:productKey", ph.Update)
	r.DELETE("/products/:id", ph.Delete)
	r.POST("/products/:id/restore", ph.Restore)

	r.GET("/ota/packages", oh.ListOTA)
	r.GET("/ota/packages/:uuid", oh.GetOTA)
	r.POST("/ota/packages", oh.CreateOTA)
	r.PUT("/ota/packages/:uuid", oh.UpdateOTA)
	r.DELETE("/ota/packages/:uuid", oh.DeleteOTA)
	r.GET("/ota/packages/:uuid/stats", oh.OTAStats)
	r.GET("/ota/packages/:uuid/batches", oh.OTABatches)
	r.GET("/ota/packages/:uuid/deployments", oh.OTADeployments)

	r.GET("/device-events", eh.ListDeviceEvents)

	return cr
}

func (cr *coverageRouters) do(method, path string, body any) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	cr.router.ServeHTTP(w, req)
	return w
}

func (cr *coverageRouters) doRaw(method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	cr.router.ServeHTTP(w, req)
	return w
}

// --- Device Handler Tests ---

func TestCoverage_Device_CreateDevice_InvalidJSON(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.doRaw("POST", "/devices", []byte("bad"))
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_CreateDevice_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().CreateDevice(gomock.Any(), gomock.Any()).Return(&model.Device{ID: 1, Name: "Test"}, nil)
	w := cr.do("POST", "/devices", map[string]any{"name": "Test", "productId": 1})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Device_CreateDevice_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().CreateDevice(gomock.Any(), gomock.Any()).Return(nil, errors.New("name is required"))
	w := cr.do("POST", "/devices", map[string]any{"name": "Test", "productId": 1})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCoverage_Device_ListDevices_SearchAndStates(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().ListDevices(gomock.Any(), 1, 10, int64(0), []string{"online"}, gomock.Any(), "test").Return([]model.Device{}, int64(0), nil)
	w := cr.do("GET", "/devices?searchText=test&states=online", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Device_ListDevices_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().ListDevices(gomock.Any(), 1, 10, int64(0), gomock.Any(), gomock.Any(), "").Return(nil, int64(0), errors.New("db error"))
	w := cr.do("GET", "/devices", nil)
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_PostTags_InvalidJSON(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.doRaw("POST", "/devices/1/tags", []byte("bad"))
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_PostTags_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().SetTags(gomock.Any(), "1", map[string]string{"k": "v"}, false).Return([]model.DeviceTag{}, nil)
	w := cr.do("POST", "/devices/1/tags", map[string]any{"tags": map[string]string{"k": "v"}})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Device_PostTags_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().SetTags(gomock.Any(), "1", map[string]string{"k": "v"}, false).Return(nil, errors.New("invalid tag key"))
	w := cr.do("POST", "/devices/1/tags", map[string]any{"tags": map[string]string{"k": "v"}})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCoverage_Device_Desired_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	desired := map[string]any{"key": "val"}
	shadow := &model.DeviceShadow{Desired: `{"key":"val"}`, Reported: `{}`, Metadata: `{}`, Version: 1}
	cr.device.EXPECT().MutateShadow(gomock.Any(), "1", int64(0), "app", &desired, nil, false).Return(shadow, nil)
	w := cr.do("PUT", "/devices/1/shadow/desired", map[string]any{"version": 0, "desired": desired})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Device_Desired_InvalidJSON(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.doRaw("PUT", "/devices/1/shadow/desired", []byte("bad"))
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_Desired_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().MutateShadow(gomock.Any(), "1", int64(0), "app", gomock.Any(), nil, false).Return(nil, repository.ErrNotFound)
	w := cr.do("PUT", "/devices/1/shadow/desired", map[string]any{"version": 0, "desired": map[string]any{"k": "v"}})
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_Reported_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	shadow := &model.DeviceShadow{Desired: `{}`, Reported: `{"temp":25}`, Metadata: `{}`, Version: 1}
	cr.device.EXPECT().MutateShadow(gomock.Any(), "1", int64(0), "device", nil, gomock.Any(), false).Return(shadow, nil)
	w := cr.do("PUT", "/devices/1/shadow/reported", map[string]any{"version": 0, "reported": map[string]any{"temp": 25}})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Device_Reported_InvalidJSON(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.doRaw("PUT", "/devices/1/shadow/reported", []byte("bad"))
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_Reported_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().MutateShadow(gomock.Any(), "1", int64(0), "device", nil, gomock.Any(), false).Return(nil, errors.New("version conflict"))
	w := cr.do("PUT", "/devices/1/shadow/reported", map[string]any{"version": 0, "reported": map[string]any{"k": "v"}})
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_ClearDesired_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	shadow := &model.DeviceShadow{Desired: `{}`, Reported: `{}`, Metadata: `{}`, Version: 1}
	cr.device.EXPECT().MutateShadow(gomock.Any(), "1", int64(0), "app", nil, nil, true).Return(shadow, nil)
	w := cr.do("DELETE", "/devices/1/shadow/desired", map[string]any{"version": 0})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Device_ClearDesired_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().MutateShadow(gomock.Any(), "1", int64(0), "app", nil, nil, true).Return(nil, errors.New("version conflict"))
	w := cr.do("DELETE", "/devices/1/shadow/desired", map[string]any{"version": 0})
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_History_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().ShadowHistory(gomock.Any(), "1").Return([]model.DeviceShadowHistory{}, nil)
	w := cr.do("GET", "/devices/1/shadow/history", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Device_History_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().ShadowHistory(gomock.Any(), "1").Return(nil, repository.ErrNotFound)
	w := cr.do("GET", "/devices/1/shadow/history", nil)
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_ShadowJSON_Delta(t *testing.T) {
	cr := newCoverageRouters(t)
	shadow := &model.DeviceShadow{
		Desired:  `{"temp":25," humidity":60}`,
		Reported: `{"temp":20}`,
		Metadata: `{}`,
		Version:  1,
	}
	cr.device.EXPECT().Shadow(gomock.Any(), "1").Return(shadow, nil)
	w := cr.do("GET", "/devices/1/shadow", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Device_BatchTemplate_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().BatchTemplate().Return([]byte("template"), nil)
	w := cr.do("GET", "/devices/batch-template", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Device_BatchTemplate_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().BatchTemplate().Return(nil, errors.New("gen error"))
	w := cr.do("GET", "/devices/batch-template", nil)
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_MockKafka_InvalidJSON(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.doRaw("POST", "/devices/mock-kafka", []byte("bad"))
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_MockKafka_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().MockKafka(gomock.Any(), gomock.Any(), "topic", "data").Return(errors.New("kafka not initialized"))
	w := cr.do("POST", "/devices/mock-kafka", map[string]any{"topic": "topic", "data": "data"})
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Device_MockKafka_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().MockKafka(gomock.Any(), gomock.Any(), "topic", "data").Return(nil)
	w := cr.do("POST", "/devices/mock-kafka", map[string]any{"topic": "topic", "data": "data"})
	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Product Handler Tests ---

func TestCoverage_Product_Create_InvalidJSON(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.doRaw("POST", "/products", []byte("bad"))
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Product_Create_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.Product{ID: 1, Name: "Test"}, nil)
	w := cr.do("POST", "/products", map[string]any{"name": "Test"})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Product_Create_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("name is required"))
	w := cr.do("POST", "/products", map[string]any{"name": "Test"})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCoverage_Product_Get_InvalidID(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.do("GET", "/products/abc", nil)
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Product_Get_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Get(gomock.Any(), int64(1)).Return(&model.Product{ID: 1, Name: "Test"}, nil)
	w := cr.do("GET", "/products/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Product_Get_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, repository.ErrNotFound)
	w := cr.do("GET", "/products/1", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCoverage_Product_List_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().List(gomock.Any(), 1, 10, "", "", "").Return(nil, int64(0), errors.New("db error"))
	w := cr.do("GET", "/products", nil)
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Product_List_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().List(gomock.Any(), 1, 10, "", "", "").Return([]model.Product{{ID: 1, Name: "Test"}}, int64(1), nil)
	w := cr.do("GET", "/products", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Product_List_WithFilters(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().List(gomock.Any(), 1, 10, "iot", "active", "search").Return([]model.Product{}, int64(0), nil)
	w := cr.do("GET", "/products?category=iot&status=active&searchText=search", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Product_Update_GetFails(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().GetByKey(gomock.Any(), "NONEXIST").Return(nil, repository.ErrNotFound)
	w := cr.do("PUT", "/products/key/NONEXIST", map[string]any{"name": "Updated"})
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCoverage_Product_Update_InvalidJSON(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().GetByKey(gomock.Any(), "P001").Return(&model.Product{ID: 1, ProductKey: "P001"}, nil)
	w := cr.doRaw("PUT", "/products/key/P001", []byte("bad"))
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Product_Update_SaveFails(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().GetByKey(gomock.Any(), "P001").Return(&model.Product{ID: 1, ProductKey: "P001"}, nil)
	cr.product.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("save error"))
	w := cr.do("PUT", "/products/key/P001", map[string]any{"name": "Updated"})
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Product_Delete_InvalidID(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.do("DELETE", "/products/abc", nil)
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Product_Delete_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Delete(gomock.Any(), int64(1)).Return(errors.New("name is required"))
	w := cr.do("DELETE", "/products/1", nil)
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Product_Delete_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)
	w := cr.do("DELETE", "/products/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Product_Restore_InvalidID(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.do("POST", "/products/abc/restore", nil)
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_Product_Restore_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Restore(gomock.Any(), int64(1)).Return(nil, repository.ErrNotFound)
	w := cr.do("POST", "/products/1/restore", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCoverage_Product_Restore_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Restore(gomock.Any(), int64(1)).Return(&model.Product{ID: 1}, nil)
	w := cr.do("POST", "/products/1/restore", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// --- OTA Handler Tests ---

func TestCoverage_OTA_List_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.ota.EXPECT().List(gomock.Any(), 1, 20).Return(nil, int64(0), errors.New("db error"))
	w := cr.do("GET", "/ota/packages", nil)
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_OTA_List_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.ota.EXPECT().List(gomock.Any(), 1, 20).Return([]model.OTAPackage{{ID: 1, PackageName: "fw"}}, int64(1), nil)
	w := cr.do("GET", "/ota/packages", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_OTA_Create_InvalidJSON(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.doRaw("POST", "/ota/packages", []byte("bad"))
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_OTA_Create_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.ota.EXPECT().Create(gomock.Any(), gomock.Any(), "P001").Return(nil)
	w := cr.do("POST", "/ota/packages", map[string]any{"packageName": "fw", "version": "1.0", "product_key": "P001"})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_OTA_Create_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.ota.EXPECT().Create(gomock.Any(), gomock.Any(), "NONEXIST").Return(repository.ErrNotFound)
	w := cr.do("POST", "/ota/packages", map[string]any{"packageName": "fw", "version": "1.0", "product_key": "NONEXIST"})
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCoverage_OTA_Update_InvalidID(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.ota.EXPECT().Get(gomock.Any(), "abc").Return(nil, repository.ErrNotFound)
	w := cr.do("PUT", "/ota/packages/abc", map[string]any{"packageName": "fw"})
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_OTA_Update_InvalidJSON(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.doRaw("PUT", "/ota/packages/1", []byte("bad"))
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_OTA_Update_GetFails(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.ota.EXPECT().Get(gomock.Any(), "1").Return(nil, repository.ErrNotFound)
	w := cr.do("PUT", "/ota/packages/1", map[string]any{"packageName": "fw"})
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCoverage_OTA_Update_SaveFails(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.ota.EXPECT().Get(gomock.Any(), "1").Return(&model.OTAPackage{ID: 1}, nil)
	cr.ota.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("save error"))
	w := cr.do("PUT", "/ota/packages/1", map[string]any{"packageName": "fw"})
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_OTA_Stats_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.ota.EXPECT().Statistics(gomock.Any(), "1").Return(service.UpgradeStatistics{PackageID: "1", TotalTargetDevices: 10}, nil)
	w := cr.do("GET", "/ota/packages/1/stats", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_OTA_Stats_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.ota.EXPECT().Statistics(gomock.Any(), "NONEXIST").Return(service.UpgradeStatistics{}, repository.ErrNotFound)
	w := cr.do("GET", "/ota/packages/NONEXIST/stats", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCoverage_OTA_Deployments_WithStatus(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.ota.EXPECT().Deployments(gomock.Any(), "1", 1, 100, "success").Return([]model.DeviceDeployment{}, int64(0), nil)
	w := cr.do("GET", "/ota/packages/1/deployments?status=success", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Device Event Handler Tests ---

func TestCoverage_DeviceEvent_List_Error(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.deviceEvent.EXPECT().List(gomock.Any(), 1, 20, "", "", "", gomock.Any(), gomock.Any()).Return(nil, int64(0), errors.New("db error"))
	w := cr.do("GET", "/device-events", nil)
	assert.True(t, w.Code >= 400, "expected error status code")
}

func TestCoverage_DeviceEvent_List_Success(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.deviceEvent.EXPECT().List(gomock.Any(), 1, 20, "temp", "D001", "temperature", gomock.Any(), gomock.Any()).Return([]model.DeviceEvent{{ID: 1, EventType: "temperature"}}, int64(1), nil)
	w := cr.do("GET", "/device-events?keyword=temp&device_key=D001&event_type=temperature", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_DeviceEvent_List_InvalidStartAt(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.do("GET", "/device-events?start_at=bad", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCoverage_DeviceEvent_List_InvalidEndAt(t *testing.T) {
	cr := newCoverageRouters(t)
	w := cr.do("GET", "/device-events?end_at=bad", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCoverage_DeviceEvent_List_WithTimeRange(t *testing.T) {
	cr := newCoverageRouters(t)
	start := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	end := time.Now().UTC().Format(time.RFC3339)
	cr.deviceEvent.EXPECT().List(gomock.Any(), 1, 20, "", "", "", gomock.Any(), gomock.Any()).Return([]model.DeviceEvent{}, int64(0), nil)
	w := cr.do("GET", fmt.Sprintf("/device-events?start_at=%s&end_at=%s", start, end), nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Response Helper Tests (page, raw, deletedAt) ---

func TestCoverage_Page_NegativePageNumber(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().List(gomock.Any(), 1, 10, "", "", "").Return([]model.Product{}, int64(0), nil)
	w := cr.do("GET", "/products?page=-1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Page_ZeroPageSize(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().List(gomock.Any(), 1, 10, "", "", "").Return([]model.Product{}, int64(0), nil)
	w := cr.do("GET", "/products?pageSize=0", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Page_NegativePageSize(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().List(gomock.Any(), 1, 10, "", "", "").Return([]model.Product{}, int64(0), nil)
	w := cr.do("GET", "/products?pageSize=-5", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Page_LargePageSize(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().List(gomock.Any(), 1, 100, "", "", "").Return([]model.Product{}, int64(0), nil)
	w := cr.do("GET", "/products?pageSize=200", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Page_NonNumericPage(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().List(gomock.Any(), 1, 10, "", "", "").Return([]model.Product{}, int64(0), nil)
	w := cr.do("GET", "/products?page=abc&pageSize=xyz", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Raw_EmptyPayload(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Get(gomock.Any(), int64(1)).Return(&model.Product{ID: 1, Metadata: ``}, nil)
	w := cr.do("GET", "/products/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Raw_LegacyString(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Get(gomock.Any(), int64(1)).Return(&model.Product{ID: 1, Metadata: `"legacy string"`}, nil)
	w := cr.do("GET", "/products/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Raw_JSONObject(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Get(gomock.Any(), int64(1)).Return(&model.Product{ID: 1, Metadata: `{"key":"value"}`}, nil)
	w := cr.do("GET", "/products/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_Raw_InvalidLegacyString(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.product.EXPECT().Get(gomock.Any(), int64(1)).Return(&model.Product{ID: 1, Metadata: `"not-valid-json`}, nil)
	w := cr.do("GET", "/products/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_DeletedAt_Valid(t *testing.T) {
	cr := newCoverageRouters(t)
	now := time.Now()
	product := &model.Product{ID: 1, Name: "Deleted", DeletedAt: gorm.DeletedAt{Time: now, Valid: true}}
	cr.product.EXPECT().Get(gomock.Any(), int64(1)).Return(product, nil)
	w := cr.do("GET", "/products/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_DeletedAt_Invalid(t *testing.T) {
	cr := newCoverageRouters(t)
	product := &model.Product{ID: 1, Name: "Active", DeletedAt: gorm.DeletedAt{Valid: false}}
	cr.product.EXPECT().Get(gomock.Any(), int64(1)).Return(product, nil)
	w := cr.do("GET", "/products/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCoverage_DeviceError_ErrVersionConflict(t *testing.T) {
	cr := newCoverageRouters(t)
	cr.device.EXPECT().MutateShadow(gomock.Any(), "1", int64(999), "app", gomock.Any(), nil, false).Return(nil, repository.ErrVersionConflict)
	w := cr.do("PUT", "/devices/1/shadow/desired", map[string]any{"version": 999, "desired": map[string]any{"k": "v"}})
	assert.Equal(t, http.StatusConflict, w.Code)
}
