package server

import (
	"context"
	"testing"

	"0things/pkg/event"
	"transport-coap/pkg/log"

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

func TestNewCoAPServer(t *testing.T) {
	v := viper.New()
	v.Set("coap.addr", ":5683")
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	mockProd := &mockEventProducer{}
	srv := NewCoAPServer(v, logger, mockProd)
	if srv == nil {
		t.Fatal("expected non-nil CoAPServer")
	}
	if srv.addr != ":5683" {
		t.Errorf("expected addr :5683, got %s", srv.addr)
	}
}
