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

func TestTelemetryHandler_QueryHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	logger := log.NewLog(viper.New())
	h := NewTelemetryHandler(&mockTelemetryService{}, logger)

	r.GET("/v1/devices/:deviceKey/telemetry/history", h.QueryHistory)

	// Test historical telemetry data.
	req, _ := http.NewRequest(http.MethodGet, "/v1/devices/dev_01/telemetry/history?property=temperature", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

}
