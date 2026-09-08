package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"0things/pkg/event"
	"transport-mqtt/pkg/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

type mockToken struct {
	err error
}

func (t *mockToken) Wait() bool                        { return true }
func (t *mockToken) WaitTimeout(_ time.Duration) bool { return true }
func (t *mockToken) Done() <-chan struct{}             { ch := make(chan struct{}); close(ch); return ch }
func (t *mockToken) Error() error                      { return t.err }

type mockMQTTClient struct {
	mqtt.Client
	publishedTopic   string
	publishedPayload []byte
	connected        bool
	publishErr       error
}

func (m *mockMQTTClient) IsConnected() bool { return m.connected }
func (m *mockMQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	m.publishedTopic = topic
	if b, ok := payload.([]byte); ok {
		m.publishedPayload = b
	}
	return &mockToken{err: m.publishErr}
}

func TestOTACommandConsumer_HandleUpgradeCommand_Success(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	mockClient := &mockMQTTClient{connected: true}
	c := NewOTACommandConsumer(mockClient, logger)

	cmd := &event.OTAUpgradeCommand{
		Transport:     "mqtt",
		BatchID:       "batch-999",
		ProductKey:    "prod_demo",
		DeviceKey:     "dev_demo_01",
		DeviceName:    "sensor_kitchen",
		TargetVersion: "v2.0.0",
		DownloadURL:   "http://example.com/ota.bin",
	}

	err := c.HandleUpgradeCommand(context.Background(), cmd, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedTopic := "/ota/device/upgrade/prod_demo/sensor_kitchen"
	if mockClient.publishedTopic != expectedTopic {
		t.Errorf("expected topic %s, got %s", expectedTopic, mockClient.publishedTopic)
	}

	var parsed event.OTAUpgradeCommand
	if err := json.Unmarshal(mockClient.publishedPayload, &parsed); err != nil {
		t.Fatalf("failed to unmarshal published payload: %v", err)
	}
	if parsed.BatchID != "batch-999" || parsed.TargetVersion != "v2.0.0" {
		t.Errorf("unexpected parsed payload: %+v", parsed)
	}
}

func TestOTACommandConsumer_HandleUpgradeCommand_DeviceNameFallback(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	mockClient := &mockMQTTClient{connected: true}
	c := NewOTACommandConsumer(mockClient, logger)

	cmd := &event.OTAUpgradeCommand{
		Transport:     "mqtt",
		BatchID:       "batch-100",
		ProductKey:    "prod_demo",
		DeviceKey:     "dev_demo_02",
		TargetVersion: "v1.1.0",
	}

	err := c.HandleUpgradeCommand(context.Background(), cmd, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedTopic := "/ota/device/upgrade/prod_demo/dev_demo_02"
	if mockClient.publishedTopic != expectedTopic {
		t.Errorf("expected topic %s, got %s", expectedTopic, mockClient.publishedTopic)
	}
}

func TestOTACommandConsumer_HandleUpgradeCommand_PublishError(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	mockClient := &mockMQTTClient{
		connected:  true,
		publishErr: errors.New("network drop"),
	}
	c := NewOTACommandConsumer(mockClient, logger)

	cmd := &event.OTAUpgradeCommand{
		Transport:     "mqtt",
		BatchID:       "batch-01",
		ProductKey:    "prod_demo",
		DeviceKey:     "dev_01",
		TargetVersion: "v1.0.0",
	}

	err := c.HandleUpgradeCommand(context.Background(), cmd, nil)
	if err == nil {
		t.Fatalf("expected error from publish, got nil")
	}
}

func TestOTACommandConsumer_HandleUpgradeCommand_DisconnectedClient(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	c := NewOTACommandConsumer(nil, logger)

	cmd := &event.OTAUpgradeCommand{
		Transport:     "mqtt",
		BatchID:       "b-01",
		ProductKey:    "prod_demo",
		DeviceKey:     "dev_demo_01",
		TargetVersion: "v1.0.1",
	}

	err := c.HandleUpgradeCommand(context.Background(), cmd, nil)
	if err == nil {
		t.Fatalf("expected error for nil/disconnected client, got nil")
	}
}
