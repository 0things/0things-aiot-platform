package server

import (
	"bytes"
	"context"
	"io"
	"testing"

	"0things/pkg/event"
	"transport-coap/pkg/log"

	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	"github.com/plgd-dev/go-coap/v3/mux"
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

type mockResponseWriter struct {
	code          codes.Code
	contentFormat message.MediaType
	body          []byte
	msg           *pool.Message
}

func (m *mockResponseWriter) Conn() mux.Conn { return nil }
func (m *mockResponseWriter) Message() *pool.Message { return m.msg }
func (m *mockResponseWriter) SetMessage(msg *pool.Message) { m.msg = msg }
func (m *mockResponseWriter) SetResponse(code codes.Code, contentFormat message.MediaType, d io.ReadSeeker, opts ...message.Option) error {
	m.code = code
	m.contentFormat = contentFormat
	if d != nil {
		b, _ := io.ReadAll(d)
		m.body = b
	}
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

func TestHandleDeviceIngress(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)
	mockProd := &mockEventProducer{}
	srv := NewCoAPServer(v, logger, mockProd)

	// 1. Missing deviceKey
	w := &mockResponseWriter{}
	poolMsg := pool.NewMessage(context.Background())
	poolMsg.SetBody(bytes.NewReader([]byte(`{"temp":25}`)))
	req := &mux.Message{
		RouteParams: &mux.RouteParams{Vars: map[string]string{}},
		Message:     poolMsg,
	}
	srv.handleDeviceIngress(w, req)
	if w.code != codes.BadRequest {
		t.Errorf("expected BadRequest, got %v", w.code)
	}

	// 2. Valid device telemetry
	w = &mockResponseWriter{}
	req.RouteParams.Vars["deviceKey"] = "coap_dev_01"
	srv.handleDeviceIngress(w, req)
	if w.code != codes.Changed {
		t.Errorf("expected Changed (2.04), got %v", w.code)
	}
	if len(mockProd.published) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(mockProd.published))
	}
	if mockProd.topics[0] != event.TopicDeviceTelemetryReport {
		t.Errorf("expected topic %s, got %s", event.TopicDeviceTelemetryReport, mockProd.topics[0])
	}
}

func TestHandleOTAProgress(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)
	mockProd := &mockEventProducer{}
	srv := NewCoAPServer(v, logger, mockProd)

	// 1. Missing deviceKey
	w := &mockResponseWriter{}
	poolMsg := pool.NewMessage(context.Background())
	poolMsg.SetBody(bytes.NewReader([]byte(`{"batch_id":"b-1","progress":50}`)))
	req := &mux.Message{
		RouteParams: &mux.RouteParams{Vars: map[string]string{}},
		Message:     poolMsg,
	}
	srv.handleOTAProgress(w, req)
	if w.code != codes.BadRequest {
		t.Errorf("expected BadRequest, got %v", w.code)
	}

	// 2. Invalid json payload
	w = &mockResponseWriter{}
	req.RouteParams.Vars["deviceKey"] = "coap_dev_01"
	poolMsgInvalid := pool.NewMessage(context.Background())
	poolMsgInvalid.SetBody(bytes.NewReader([]byte(`not a json`)))
	req.Message = poolMsgInvalid
	srv.handleOTAProgress(w, req)
	if w.code != codes.BadRequest {
		t.Errorf("expected BadRequest for invalid json, got %v", w.code)
	}

	// 3. Valid OTA progress report
	w = &mockResponseWriter{}
	req.Message = poolMsg
	srv.handleOTAProgress(w, req)
	if w.code != codes.Changed {
		t.Errorf("expected Changed, got %v", w.code)
	}
	if len(mockProd.published) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(mockProd.published))
	}
	if mockProd.topics[0] != event.TopicOTAProgressReport {
		t.Errorf("expected topic %s, got %s", event.TopicOTAProgressReport, mockProd.topics[0])
	}
	rep, ok := mockProd.published[0].(*event.OTAUpgradeReport)
	if !ok || rep.DeviceKey != "coap_dev_01" || rep.BatchID != "b-1" {
		t.Errorf("unexpected report: %+v", rep)
	}
}
