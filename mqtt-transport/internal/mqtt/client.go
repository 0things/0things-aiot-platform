package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
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
func NewService(config *viper.Viper, logger *zap.Logger, producer *kafka.Producer) *Service {
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

	if username := config.GetString("mqtt.username"); username != "" {
		opts.SetUsername(username)
	}
	if password := config.GetString("mqtt.password"); password != "" {
		opts.SetPassword(password)
	}

	svc.client = mqtt.NewClient(opts)
	return svc
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

// dispatchUplink 统一提取上行报文、解析 deviceKey 并投递到 Kafka 专属 Topic
func (s *Service) dispatchUplink(msgType string, topic string, payload []byte) {
	deviceKey := ExtractDeviceKey(topic)
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from topic", zap.String("topic", topic), zap.String("msg_type", msgType))
		return
	}

	deviceMsg := model.DeviceMessage{
		DeviceKey:   deviceKey,
		Transport:   "mqtt",
		MessageType: msgType,
		Payload:     json.RawMessage(payload),
		Timestamp:   time.Now().UTC(),
		Headers:     map[string]string{"topic": topic},
	}

	if err := s.producer.SendDeviceMessage(context.Background(), deviceMsg); err != nil {
		s.logger.Error("failed to publish uplink message to kafka",
			zap.String("topic", topic),
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

// handleOtaProgress 专职处理 OTA 固件升级进度 ➔ 投递至 ota.report.v1
func (s *Service) handleOtaProgress(_ mqtt.Client, msg mqtt.Message) {
	s.dispatchUplink("ota_report", msg.Topic(), msg.Payload())
}

// handleDeviceEvent 专职处理设备特定告警与事件 ➔ 投递至 device.event.v1
func (s *Service) handleDeviceEvent(_ mqtt.Client, msg mqtt.Message) {
	// 避免与 handleTelemetry 重复：若为 /thing/event/property/post 则由 handleTelemetry 处理，此处跳过
	if strings.HasSuffix(msg.Topic(), "/thing/event/property/post") {
		return
	}
	s.dispatchUplink("event", msg.Topic(), msg.Payload())
}

// HandleDownlinkCommand 处理来自 Kafka 的下行控制指令并推向设备（使用标准枚举模板）
func (s *Service) HandleDownlinkCommand(ctx context.Context, cmd model.DeviceCommand) error {
	// 校验 deviceKey 与 endpoint 是否包含非法 MQTT 通配符防止注入
	if strings.ContainsAny(cmd.DeviceKey, "+#\n\r") {
		return fmt.Errorf("invalid device key containing MQTT wildcard: %s", cmd.DeviceKey)
	}

	topic := cmd.Endpoint
	if topic == "" {
		switch cmd.CommandType {
		case "ota_upgrade", "ota":
			topic = fmt.Sprintf(enum.MQTTTplOTAUpgrade, cmd.DeviceKey)
		default:
			topic = fmt.Sprintf(enum.MQTTTplPropertySet, cmd.DeviceKey)
		}
	} else if strings.ContainsAny(topic, "+#\n\r") {
		return fmt.Errorf("invalid custom endpoint topic containing MQTT wildcard: %s", topic)
	}

	token := s.client.Publish(topic, 1, false, cmd.Payload)
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("publish to device timeout: %s", topic)
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("publish to device failed: %w", err)
	}

	s.logger.Info("downlink command published to device",
		zap.String("topic", topic),
		zap.String("device_key", cmd.DeviceKey),
		zap.String("command_type", cmd.CommandType),
	)
	return nil
}

// ExtractDeviceKey 统一从主题路径中提取设备唯一标识符 deviceKey。
// 支持格式：
// 1. 标准属性与事件上报: /sys/{deviceKey}/thing/...
// 2. OTA 进度上报: /sys/{deviceKey}/ota/...
func ExtractDeviceKey(topic string) string {
	parts := strings.Split(topic, "/")
	if len(parts) >= 3 && parts[1] == "sys" {
		return parts[2]
	}
	return ""
}
