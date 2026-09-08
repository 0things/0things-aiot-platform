package server

import (
	"context"
	"fmt"
	"testing"

	"0things/pkg/event"
	"transport-mqtt/internal/enum"
	"transport-mqtt/pkg/log"

	"github.com/spf13/viper"
)

func TestExtractDeviceKey(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		expected string
	}{
		{
			name:     "standard telemetry topic",
			topic:    "/sys/prod_xyz/device_abc_123/thing/event/property/post",
			expected: "device_abc_123",
		},
		{
			name:     "standard ota progress topic",
			topic:    "/sys/prod_xyz/sensor_999/ota/device/progress",
			expected: "sensor_999",
		},
		{
			name:     "custom event topic",
			topic:    "/sys/prod_xyz/meter_001/thing/event/high_voltage/post",
			expected: "meter_001",
		},
		{
			name:     "independent ota progress topic",
			topic:    "/ota/device/progress/prod_xyz/sensor_999",
			expected: "sensor_999",
		},
		{
			name:     "invalid topic prefix",
			topic:    "/other/prod_xyz/device_abc_123/property/post",
			expected: "",
		},
		{
			name:     "empty topic",
			topic:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ExtractDeviceKey(tt.topic)
			if actual != tt.expected {
				t.Errorf("ExtractDeviceKey(%q) = %q; want %q", tt.topic, actual, tt.expected)
			}
		})
	}
}

func TestExtractProductKey(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		expected string
	}{
		{
			name:     "standard telemetry topic",
			topic:    "/sys/prod_xyz/device_abc_123/thing/event/property/post",
			expected: "prod_xyz",
		},
		{
			name:     "standard ota progress topic",
			topic:    "/sys/prod_xyz/sensor_999/ota/device/progress",
			expected: "prod_xyz",
		},
		{
			name:     "independent ota inform topic",
			topic:    "/ota/device/inform/prod_xyz/sensor_999",
			expected: "prod_xyz",
		},
		{
			name:     "invalid topic prefix",
			topic:    "/other/prod_xyz/device_abc_123/property/post",
			expected: "",
		},
		{
			name:     "empty topic",
			topic:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ExtractProductKey(tt.topic)
			if actual != tt.expected {
				t.Errorf("ExtractProductKey(%q) = %q; want %q", tt.topic, actual, tt.expected)
			}
		})
	}
}

func TestDownlinkTopicFormatting(t *testing.T) {
	productKey := "prod_xyz"
	deviceKey := "test_dev_01"

	propertyTopic := fmt.Sprintf(enum.MQTTTplPropertySet, productKey, deviceKey)
	if propertyTopic != "/sys/prod_xyz/test_dev_01/thing/service/property/set" {
		t.Errorf("unexpected property downlink topic: %s", propertyTopic)
	}

	otaTopic := fmt.Sprintf(enum.MQTTTplOTAUpgrade, productKey, deviceKey)
	if otaTopic != "/sys/prod_xyz/test_dev_01/ota/device/upgrade" {
		t.Errorf("unexpected ota downlink topic: %s", otaTopic)
	}

	otaTopicV1 := fmt.Sprintf(enum.MQTTTplOTAUpgradeV1, productKey, deviceKey)
	if otaTopicV1 != "/ota/device/upgrade/prod_xyz/test_dev_01" {
		t.Errorf("unexpected ota v1 downlink topic: %s", otaTopicV1)
	}
}

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

type fakeMqttMessage struct {
	topic   string
	payload []byte
}

func (f *fakeMqttMessage) Duplicate() bool   { return false }
func (f *fakeMqttMessage) Qos() byte         { return 0 }
func (f *fakeMqttMessage) Retained() bool    { return false }
func (f *fakeMqttMessage) Topic() string     { return f.topic }
func (f *fakeMqttMessage) MessageID() uint16 { return 1 }
func (f *fakeMqttMessage) Payload() []byte   { return f.payload }
func (f *fakeMqttMessage) Ack()              {}

func TestHandleOtaProgressPublishEvent(t *testing.T) {
	mockProd := &mockEventProducer{}
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	svc := &MQTTServer{
		logger:        logger,
		eventProducer: mockProd,
	}

	rawPayload := []byte(`{"batch_id":"b-123","progress":50,"stage":"DOWNLOADING"}`)
	msg := &fakeMqttMessage{
		topic:   "/sys/prod_demo/dev_demo_01/ota/device/progress",
		payload: rawPayload,
	}

	svc.handleOtaProgress(nil, msg)

	if len(mockProd.published) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(mockProd.published))
	}
	if mockProd.topics[0] != event.TopicOTAProgressReport {
		t.Errorf("expected topic %s, got %s", event.TopicOTAProgressReport, mockProd.topics[0])
	}
	report, ok := mockProd.published[0].(*event.OTAUpgradeReport)
	if !ok {
		t.Fatalf("expected *event.OTAUpgradeReport, got %T", mockProd.published[0])
	}
	if report.BatchID != "b-123" || report.DeviceKey != "dev_demo_01" || report.ProductKey != "prod_demo" {
		t.Errorf("unexpected report: %+v", report)
	}
}

func TestHandleTelemetryPublishEvent(t *testing.T) {
	mockProd := &mockEventProducer{}
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	svc := &MQTTServer{
		logger:        logger,
		eventProducer: mockProd,
	}

	rawPayload := []byte(`{"temperature": 25.5, "humidity": 60}`)
	msg := &fakeMqttMessage{
		topic:   "/sys/prod_demo/dev_demo_01/thing/event/property/post",
		payload: rawPayload,
	}

	svc.handleTelemetry(nil, msg)

	if len(mockProd.published) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(mockProd.published))
	}
	if mockProd.topics[0] != event.TopicDeviceTelemetryReport {
		t.Errorf("expected topic %s, got %s", event.TopicDeviceTelemetryReport, mockProd.topics[0])
	}
	deviceMsg, ok := mockProd.published[0].(*event.DeviceMessage)
	if !ok {
		t.Fatalf("expected *event.DeviceMessage, got %T", mockProd.published[0])
	}
	if deviceMsg.DeviceKey != "dev_demo_01" || deviceMsg.ProductKey != "prod_demo" || deviceMsg.MessageType != "telemetry" {
		t.Errorf("unexpected deviceMsg: %+v", deviceMsg)
	}
}

func TestHandleDeviceEventPublishEvent(t *testing.T) {
	mockProd := &mockEventProducer{}
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	svc := &MQTTServer{
		logger:        logger,
		eventProducer: mockProd,
	}

	rawPayload := []byte(`{"error_code": 1001, "msg": "overheat"}`)
	msg := &fakeMqttMessage{
		topic:   "/sys/prod_demo/dev_demo_01/thing/event/overheat_alarm/post",
		payload: rawPayload,
	}

	svc.handleDeviceEvent(nil, msg)

	if len(mockProd.published) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(mockProd.published))
	}
	if mockProd.topics[0] != event.TopicDeviceEventReport {
		t.Errorf("expected topic %s, got %s", event.TopicDeviceEventReport, mockProd.topics[0])
	}
	deviceMsg, ok := mockProd.published[0].(*event.DeviceMessage)
	if !ok {
		t.Fatalf("expected *event.DeviceMessage, got %T", mockProd.published[0])
	}
	if deviceMsg.DeviceKey != "dev_demo_01" || deviceMsg.ProductKey != "prod_demo" || deviceMsg.MessageType != "event" {
		t.Errorf("unexpected deviceMsg: %+v", deviceMsg)
	}
}
