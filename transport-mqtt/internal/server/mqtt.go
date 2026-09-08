package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"transport-mqtt/internal/enum"
	"transport-mqtt/pkg/log"
	"transport-mqtt/pkg/server"

	"0things/pkg/event"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// MQTTServer manages MQTT broker connections, subscriptions and event publishing.
type MQTTServer struct {
	client        mqtt.Client
	logger        *log.Logger
	eventProducer event.Producer
	mu            sync.Mutex
}

var _ server.Server = (*MQTTServer)(nil)

// NewMQTTServer initializes the MQTT transport server.
func NewMQTTServer(config *viper.Viper, logger *log.Logger, eventProducer event.Producer) (*MQTTServer, error) {
	broker := config.GetString("mqtt.broker")
	if broker == "" {
		broker = "tcp://127.0.0.1:1883"
	}

	clientID := config.GetString("mqtt.client_id")
	if clientID == "" {
		clientID = "0things-mqtt-transport-" + uuid.NewString()[:8]
	}

	svc := &MQTTServer{
		logger:        logger,
		eventProducer: eventProducer,
	}

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetKeepAlive(60 * time.Second).
		SetOnConnectHandler(func(c mqtt.Client) {
			svc.logger.Info("connected/reconnected to MQTT broker, registering subscriptions...")
			svc.registerSubscriptions(c)
		}).
		SetConnectionLostHandler(func(_ mqtt.Client, err error) {
			svc.logger.Warn("MQTT connection lost, waiting for auto-reconnect", zap.Error(err))
		})

	tlsConfig, err := mqttTLSConfig(config)
	if err != nil {
		return nil, err
	}
	if tlsConfig != nil {
		opts.SetTLSConfig(tlsConfig)
	}

	if username := config.GetString("mqtt.username"); username != "" {
		opts.SetUsername(username)
	}
	if password := config.GetString("mqtt.password"); password != "" {
		opts.SetPassword(password)
	}

	svc.client = mqtt.NewClient(opts)
	return svc, nil
}

func mqttTLSConfig(config *viper.Viper) (*tls.Config, error) {
	caFile := config.GetString("mqtt.tls.ca_file")
	certFile := config.GetString("mqtt.tls.cert_file")
	keyFile := config.GetString("mqtt.tls.key_file")
	if caFile == "" && certFile == "" && keyFile == "" {
		return nil, nil
	}
	if certFile == "" || keyFile == "" {
		return nil, fmt.Errorf("mqtt TLS requires both cert_file and key_file")
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("load mqtt client certificate: %w", err)
	}
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12, Certificates: []tls.Certificate{cert}}
	if caFile != "" {
		pem, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("read mqtt CA certificate: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("parse mqtt CA certificate")
		}
		tlsConfig.RootCAs = pool
	}
	return tlsConfig, nil
}

// Start connects to the broker and keeps running until context is cancelled.
func (s *MQTTServer) Start(ctx context.Context) error {
	s.mu.Lock()
	token := s.client.Connect()
	s.mu.Unlock()

	if !token.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("mqtt broker connect timeout")
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt broker connect failed: %w", err)
	}

	s.logger.Info("MQTT server connected and listening...")
	<-ctx.Done()
	return nil
}

// Stop gracefully stops the MQTT server and disconnects from broker.
func (s *MQTTServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client != nil && s.client.IsConnected() {
		s.client.Disconnect(250)
	}
	s.logger.Info("MQTT server stopped")
	return nil
}

// registerSubscriptions registers handlers for incoming MQTT topics.
func (s *MQTTServer) registerSubscriptions(c mqtt.Client) {
	s.subscribeTopicWithClient(c, enum.MQTTSubTelemetry, s.handleTelemetry)
	s.subscribeTopicWithClient(c, enum.MQTTSubOTAProgress, s.handleOtaProgress)
	s.subscribeTopicWithClient(c, enum.MQTTSubOTAProgressV1, s.handleOtaProgress)
	s.subscribeTopicWithClient(c, enum.MQTTSubOTAInform, s.handleOtaProgress)
	s.subscribeTopicWithClient(c, enum.MQTTSubEvent, s.handleDeviceEvent)
}

func (s *MQTTServer) subscribeTopicWithClient(c mqtt.Client, topic string, handler mqtt.MessageHandler) {
	subToken := c.Subscribe(topic, 0, handler)
	if subToken.Wait() && subToken.Error() != nil {
		s.logger.Error("failed to subscribe to topic", zap.String("topic", topic), zap.Error(subToken.Error()))
	} else {
		s.logger.Info("registered dedicated handler for MQTT topic", zap.String("topic", topic))
	}
}

func (s *MQTTServer) handleTelemetry(_ mqtt.Client, msg mqtt.Message) {
	deviceKey := ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from telemetry topic", zap.String("topic", msg.Topic()))
		return
	}
	productKey := ExtractProductKey(msg.Topic())

	deviceMsg := event.DeviceMessage{
		DeviceKey:   deviceKey,
		ProductKey:  productKey,
		Transport:   "mqtt",
		MessageType: "telemetry",
		Payload:     json.RawMessage(msg.Payload()),
		Timestamp:   time.Now().UTC(),
		Headers:     map[string]string{"topic": msg.Topic()},
	}

	s.logger.Info("received MQTT telemetry message",
		zap.String("topic", msg.Topic()),
		zap.String("device_key", deviceKey),
		zap.String("product_key", productKey),
		zap.Int("payload_bytes", len(msg.Payload())),
	)

	if s.eventProducer != nil {
		if err := s.eventProducer.Publish(context.Background(), event.TopicDeviceTelemetryReport, &deviceMsg); err != nil {
			s.logger.Error("failed to publish device telemetry event", zap.Error(err))
		}
	}
}

func (s *MQTTServer) handleDeviceEvent(_ mqtt.Client, msg mqtt.Message) {
	if strings.HasSuffix(msg.Topic(), "/thing/event/property/post") {
		return
	}

	deviceKey := ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from event topic", zap.String("topic", msg.Topic()))
		return
	}
	productKey := ExtractProductKey(msg.Topic())

	deviceMsg := event.DeviceMessage{
		DeviceKey:   deviceKey,
		ProductKey:  productKey,
		Transport:   "mqtt",
		MessageType: "event",
		Payload:     json.RawMessage(msg.Payload()),
		Timestamp:   time.Now().UTC(),
		Headers:     map[string]string{"topic": msg.Topic()},
	}

	s.logger.Info("received MQTT event message",
		zap.String("topic", msg.Topic()),
		zap.String("device_key", deviceKey),
		zap.String("product_key", productKey),
		zap.Int("payload_bytes", len(msg.Payload())),
	)

	if s.eventProducer != nil {
		if err := s.eventProducer.Publish(context.Background(), event.TopicDeviceEventReport, &deviceMsg); err != nil {
			s.logger.Error("failed to publish device event report", zap.Error(err))
		}
	}
}

func (s *MQTTServer) handleOtaProgress(_ mqtt.Client, msg mqtt.Message) {
	deviceKey := ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from OTA topic", zap.String("topic", msg.Topic()))
		return
	}

	var report event.OTAUpgradeReport
	if err := json.Unmarshal(msg.Payload(), &report); err != nil {
		s.logger.Warn("invalid OTA report payload", zap.String("topic", msg.Topic()), zap.Error(err))
		return
	}
	if report.DeviceKey == "" {
		report.DeviceKey = deviceKey
	}
	if report.ProductKey == "" {
		if productKey := ExtractProductKey(msg.Topic()); productKey != "" {
			report.ProductKey = productKey
		}
	}
	if report.ReportedAt.IsZero() {
		report.ReportedAt = time.Now()
	}

	s.logger.Info("received OTA progress report",
		zap.String("topic", msg.Topic()),
		zap.String("device_key", deviceKey),
		zap.Any("report", report),
	)

	if s.eventProducer != nil {
		if err := s.eventProducer.Publish(context.Background(), event.TopicOTAProgressReport, &report); err != nil {
			s.logger.Error("failed to publish OTA progress report event", zap.Error(err))
		}
	}
}

// ExtractDeviceKey extracts deviceKey from topic path.
func ExtractDeviceKey(topic string) string {
	parts := strings.Split(topic, "/")
	if len(parts) >= 4 && parts[1] == "sys" {
		return parts[3]
	}
	if len(parts) >= 6 && parts[1] == "ota" && parts[2] == "device" && (parts[3] == "progress" || parts[3] == "inform") {
		return parts[5]
	}
	return ""
}

// ExtractProductKey extracts productKey from topic path.
func ExtractProductKey(topic string) string {
	parts := strings.Split(topic, "/")
	if len(parts) >= 3 && parts[1] == "sys" {
		return parts[2]
	}
	if len(parts) >= 6 && parts[1] == "ota" && parts[2] == "device" && (parts[3] == "progress" || parts[3] == "inform") {
		return parts[4]
	}
	return ""
}
