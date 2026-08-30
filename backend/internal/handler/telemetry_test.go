package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"aiot-backend/internal/model"
	"aiot-backend/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type mockTelemetryService struct{}

func (m *mockTelemetryService) QueryHistory(ctx context.Context, req model.TelemetryQueryReq) ([]model.TelemetryPoint, error) {
	return []model.TelemetryPoint{
		{Timestamp: time.Now().UnixMilli(), Property: req.Property, Value: 26.5},
	}, nil
}

func (m *mockTelemetryService) GetShadow(ctx context.Context, deviceKey string) (*model.DeviceShadowSnapshot, error) {
	return &model.DeviceShadowSnapshot{
		DeviceKey:  deviceKey,
		Attributes: map[string]interface{}{"temperature": 26.5},
		LastSeen:   time.Now(),
	}, nil
}

func TestTelemetryHandler_QueryHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	logger := log.NewLog(viper.New())
	h := NewTelemetryHandler(&mockTelemetryService{}, logger)

	r.GET("/v1/devices/:deviceKey/telemetry/history", h.QueryHistory)
	r.GET("/v1/devices/:deviceKey/shadow", h.GetShadow)

	// 测试历史曲线
	req, _ := http.NewRequest(http.MethodGet, "/v1/devices/dev_01/telemetry/history?property=temperature", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	// 测试设备影子
	reqShadow, _ := http.NewRequest(http.MethodGet, "/v1/devices/dev_01/shadow", nil)
	wShadow := httptest.NewRecorder()
	r.ServeHTTP(wShadow, reqShadow)

	if wShadow.Code != http.StatusOK {
		t.Fatalf("expected 200 for shadow, got %d", wShadow.Code)
	}
}
