package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"transport-mqtt/internal/consumer"
	"transport-mqtt/internal/enum"
	"transport-mqtt/internal/handler"
	"transport-mqtt/pkg/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

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

func TestMQTTTLSConfig(t *testing.T) {
	// 1. No TLS configured
	v := viper.New()
	tlsConfig, err := MQTTTLSConfig(v)
	if err != nil {
		t.Fatalf("expected nil error when no TLS configured, got %v", err)
	}
	if tlsConfig != nil {
		t.Errorf("expected nil tlsConfig, got %v", tlsConfig)
	}

	// 2. Only cert_file configured (missing key_file)
	v.Set("mqtt.tls.cert_file", "cert.pem")
	_, err = MQTTTLSConfig(v)
	if err == nil {
		t.Fatalf("expected error when key_file is missing, got nil")
	}
}

type mockToken struct {
	err error
}

func (t *mockToken) Wait() bool                        { return true }
func (t *mockToken) WaitTimeout(_ time.Duration) bool { return true }
func (t *mockToken) Done() <-chan struct{}             { ch := make(chan struct{}); close(ch); return ch }
func (t *mockToken) Error() error                      { return t.err }

type mockClientForServer struct {
	mqtt.Client
	subscribedTopics []string
	connected        bool
	disconnected     bool
}

func (m *mockClientForServer) Connect() mqtt.Token {
	m.connected = true
	return &mockToken{}
}

func (m *mockClientForServer) Disconnect(quiesce uint) {
	m.disconnected = true
	m.connected = false
}

func (m *mockClientForServer) IsConnected() bool {
	return m.connected
}

func (m *mockClientForServer) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	m.subscribedTopics = append(m.subscribedTopics, topic)
	return &mockToken{}
}

func TestRegisterSubscriptions(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	mockClient := &mockClientForServer{}
	ingressHandler := handler.NewIngressHandler(nil, logger)

	// 1. Nil handler does nothing
	RegisterSubscriptions(mockClient, nil, logger)
	if len(mockClient.subscribedTopics) != 0 {
		t.Fatalf("expected 0 subscriptions for nil handler, got %d", len(mockClient.subscribedTopics))
	}

	// 2. Valid handler registers all expected topics
	RegisterSubscriptions(mockClient, ingressHandler, logger)
	if len(mockClient.subscribedTopics) != 5 {
		t.Errorf("expected 5 subscriptions, got %d", len(mockClient.subscribedTopics))
	}
}

func TestNewMQTTServer_And_Stop(t *testing.T) {
	v := viper.New()
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	mockClient := &mockClientForServer{connected: true}
	consumerMgr := consumer.NewManager(nil, nil, logger)

	server := NewMQTTServer(mockClient, logger, consumerMgr)
	if server == nil {
		t.Fatalf("expected non-nil server")
	}

	err := server.Stop(context.Background())
	if err != nil {
		t.Fatalf("unexpected error on stop: %v", err)
	}
	if !mockClient.disconnected {
		t.Errorf("expected mock client to be disconnected")
	}
}

func TestNewMQTTClient_Configuration(t *testing.T) {
	v := viper.New()
	v.Set("mqtt.broker", "tcp://127.0.0.1:1883")
	v.Set("mqtt.client_id", "test-client")
	v.Set("mqtt.username", "admin")
	v.Set("mqtt.password", "secret")
	v.Set("log.mode", "console")
	v.Set("log.log_level", "error")
	logger := log.NewLog(v)

	client, err := NewMQTTClient(v, logger, nil)
	if err != nil {
		t.Fatalf("unexpected error creating mqtt client: %v", err)
	}
	if client == nil {
		t.Fatalf("expected non-nil client")
	}
}
