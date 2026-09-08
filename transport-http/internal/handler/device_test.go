package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"0things/pkg/event"
	"transport-http/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type mockEventProducer struct {
	published []any
	topics    []event.Topic
}

func (m *mockEventProducer) Publish(ctx context.Context, topic event.Topic, payload any, opts ...event.PublishOption) error {
	m.topics = append(m.topics, topic)
	m.published = append(m.published, payload)
	return nil
}

func (m *mockEventProducer) Close() error {
	return nil
}

func setupTestRouter(producer event.Producer) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)
	h := NewDeviceHandler(logger, producer)

	api := r.Group("/api/v1/:deviceKey")
	{
		api.POST("/telemetry", h.PostTelemetry)
		api.POST("/attributes", h.PostAttributes)
		api.POST("/ota/progress", h.PostOtaProgress)
		api.POST("/events/:eventType", h.PostEvent)
	}
	return r
}

func TestPostTelemetry(t *testing.T) {
	mockProd := &mockEventProducer{}
	r := setupTestRouter(mockProd)

	body := []byte(`{"temperature": 26.8, "humidity": 60}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/dev_http_01/telemetry", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d, body: %s", w.Code, w.Body.String())
	}
	if len(mockProd.published) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(mockProd.published))
	}
	if mockProd.topics[0] != event.TopicDeviceTelemetryReport {
		t.Errorf("expected topic %s, got %s", event.TopicDeviceTelemetryReport, mockProd.topics[0])
	}
}

func TestPostOtaProgress(t *testing.T) {
	mockProd := &mockEventProducer{}
	r := setupTestRouter(mockProd)

	body := []byte(`{"batch_id": "b-1", "step": 50, "stage": "DOWNLOADING"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/dev_http_01/ota/progress", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d, body: %s", w.Code, w.Body.String())
	}
	if len(mockProd.published) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(mockProd.published))
	}
	if mockProd.topics[0] != event.TopicOTAProgressReport {
		t.Errorf("expected topic %s, got %s", event.TopicOTAProgressReport, mockProd.topics[0])
	}
}

func TestPostEmptyBody(t *testing.T) {
	mockProd := &mockEventProducer{}
	r := setupTestRouter(mockProd)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/dev_http_01/telemetry", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for empty body, got %d", w.Code)
	}
}
