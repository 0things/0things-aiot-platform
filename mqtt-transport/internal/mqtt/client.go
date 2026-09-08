package mqtt

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

	"mqtt-transport/internal/enum"
	"mqtt-transport/internal/kafka"
	"mqtt-transport/internal/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Service 负责管理与 MQTT Broker 的连接，并将标准 MQTT 主题精准绑定至专属 Handler。
type Service struct {
	client   mqtt.Client
	producer *kafka.Producer
	logger   *zap.Logger
	mu       sync.Mutex
}

// NewService 初始化 MQTT 客户端。
func NewService(config *viper.Viper, logger *zap.Logger, producer *kafka.Producer) (*Service, error) {
	broker := config.GetString("mqtt.broker")
	if broker == "" {
		broker = "tcp://127.0.0.1:1883"
	}

	clientID := config.GetString("mqtt.client_id")
	if clientID == "" {
		clientID = "0things-mqtt-transport-" + uuid.NewString()[:8]
	}

	svc := &Service{
		producer: producer,
		logger:   logger,
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

// Start 建立连接，并保持服务运行。
func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	token := s.client.Connect()
	s.mu.Unlock()

	if !token.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("mqtt broker connect timeout")
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt broker connect failed: %w", err)
	}

	s.logger.Info("MQTT client connected and listening...")

	<-ctx.Done()
	s.client.Disconnect(250)
	s.logger.Info("MQTT client disconnected")
	return nil
}

// registerSubscriptions 在初次连接与重连时统一恢复主题订阅
func (s *Service) registerSubscriptions(c mqtt.Client) {
	// 1. 标准时序遥测上报 ➔ 绑定 handleTelemetry
	s.subscribeTopicWithClient(c, enum.MQTTSubTelemetry, s.handleTelemetry)

	// 2. 标准 OTA 升级进度上报 ➔ 绑定 handleOtaProgress
	s.subscribeTopicWithClient(c, enum.MQTTSubOTAProgress, s.handleOtaProgress)
	s.subscribeTopicWithClient(c, enum.MQTTSubOTAProgressV1, s.handleOtaProgress)
	s.subscribeTopicWithClient(c, enum.MQTTSubOTAInform, s.handleOtaProgress)

	// 3. 标准业务事件上报 ➔ 绑定 handleDeviceEvent
	s.subscribeTopicWithClient(c, enum.MQTTSubEvent, s.handleDeviceEvent)
}

func (s *Service) subscribeTopicWithClient(c mqtt.Client, topic string, handler mqtt.MessageHandler) {
	subToken := c.Subscribe(topic, 0, handler)
	if subToken.Wait() && subToken.Error() != nil {
		s.logger.Error("failed to subscribe to topic", zap.String("topic", topic), zap.Error(subToken.Error()))
	} else {
		s.logger.Info("registered dedicated handler for MQTT topic", zap.String("topic", topic))
	}
}

// dispatchUplink 统一提取上行报文、解析 deviceKey 与 productKey 并投递到 Kafka 专属 Topic
func (s *Service) dispatchUplink(msgType string, topic string, payload []byte) {
	deviceKey := ExtractDeviceKey(topic)
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from topic", zap.String("topic", topic), zap.String("msg_type", msgType))
		return
	}
	productKey := ExtractProductKey(topic)

	deviceMsg := model.DeviceMessage{
		DeviceKey:   deviceKey,
		ProductKey:  productKey,
		Transport:   "mqtt",
		MessageType: msgType,
		Payload:     json.RawMessage(payload),
		Timestamp:   time.Now().UTC(),
		Headers:     map[string]string{"topic": topic},
	}

	if err := s.producer.SendDeviceMessage(context.Background(), deviceMsg); err != nil {
		s.logger.Error("failed to publish uplink message to kafka",
			zap.String("topic", topic),
			zap.String("product_key", productKey),
			zap.String("device_key", deviceKey),
			zap.String("msg_type", msgType),
			zap.Error(err),
		)
	}
}

// handleTelemetry 专职处理时序遥测与属性 ➔ 投递至 device.telemetry.v1
func (s *Service) handleTelemetry(_ mqtt.Client, msg mqtt.Message) {
	s.dispatchUplink("telemetry", msg.Topic(), msg.Payload())
}

// handleOtaProgress 专职处理 OTA 固件升级进度 ➔ 投递至 ota.upgrade.report.v1
func (s *Service) handleOtaProgress(_ mqtt.Client, msg mqtt.Message) {
	deviceKey := ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from OTA topic", zap.String("topic", msg.Topic()))
		return
	}

	var report map[string]interface{}
	if err := json.Unmarshal(msg.Payload(), &report); err != nil {
		s.logger.Warn("invalid OTA report payload", zap.String("topic", msg.Topic()), zap.Error(err))
		return
	}
	// The topic authenticates the producer identity; preserve an explicit
	// payload device_key when present, otherwise fill it from the topic.
	if _, ok := report["device_key"]; !ok {
		report["device_key"] = deviceKey
	}
	if _, ok := report["product_key"]; !ok {
		if productKey := ExtractProductKey(msg.Topic()); productKey != "" {
			report["product_key"] = productKey
		}
	}
	if err := s.producer.SendOTAReport(context.Background(), deviceKey, report); err != nil {
		s.logger.Error("failed to publish OTA report to kafka", zap.String("topic", msg.Topic()), zap.String("device_key", deviceKey), zap.Error(err))
	}
}

// handleDeviceEvent 专职处理设备特定告警与事件 ➔ 投递至 device.event.v1
func (s *Service) handleDeviceEvent(_ mqtt.Client, msg mqtt.Message) {
	// 避免与 handleTelemetry 重复：若为 /thing/event/property/post 则由 handleTelemetry 处理，此处跳过
	if strings.HasSuffix(msg.Topic(), "/thing/event/property/post") {
		return
	}
	s.dispatchUplink("event", msg.Topic(), msg.Payload())
}

// ExtractDeviceKey 统一从主题路径中提取设备唯一标识符 deviceKey。
// 支持格式：
// 1. 标准物模型及OTA上报: /sys/{productKey}/{deviceKey}/...
// 2. 独立OTA路径上报: /ota/device/{action}/{productKey}/{deviceKey}
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

// ExtractProductKey 统一从主题路径中提取产品标识符 productKey。
// 支持格式：
// 1. 标准物模型及OTA上报: /sys/{productKey}/{deviceKey}/...
// 2. 独立OTA路径上报: /ota/device/{action}/{productKey}/{deviceKey}
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
