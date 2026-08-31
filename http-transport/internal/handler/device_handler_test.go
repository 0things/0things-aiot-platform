package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"http-transport/internal/kafka"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	v := viper.New()
	logger := zap.NewNop()
	producer, _, _ := kafka.NewProducer(v, logger)
	h := NewDeviceHandler(producer, logger)

	api := r.Group("/api/v1/:deviceKey")
	{
		api.POST("/telemetry", h.PostTelemetry)
		api.POST("/ota/progress", h.PostOtaProgress)
		api.POST("/events/:eventType", h.PostEvent)
	}
	return r
}

func TestPostTelemetry(t *testing.T) {
	r := setupTestRouter()

	body := []byte(`{"temperature": 26.8, "humidity": 60}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/dev_http_01/telemetry", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestPostOtaProgress(t *testing.T) {
	r := setupTestRouter()

	body := []byte(`{"step": 50, "desc": "downloading firmware"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/dev_http_01/ota/progress", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestPostEmptyBody(t *testing.T) {
	r := setupTestRouter()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/dev_http_01/telemetry", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for empty body, got %d", w.Code)
	}
}
