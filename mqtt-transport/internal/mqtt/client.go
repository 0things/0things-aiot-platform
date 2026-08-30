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

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetKeepAlive(60 * time.Second)

	if username := config.GetString("mqtt.username"); username != "" {
		opts.SetUsername(username)
	}
	if password := config.GetString("mqtt.password"); password != "" {
		opts.SetPassword(password)
	}

	return &Service{
		client:   mqtt.NewClient(opts),
		producer: producer,
		logger:   logger,
	}
}

// Start 建立连接，并直接通过标准枚举注册专属主题 Handler。
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

	s.logger.Info("connected to MQTT broker successfully")

	// 1. 标准时序遥测上报 ➔ 绑定 handleTelemetry
	s.subscribeTopic(enum.MQTTSubTelemetry, s.handleTelemetry)

	// 2. 标准 OTA 升级进度上报 ➔ 绑定 handleOtaProgress
	s.subscribeTopic(enum.MQTTSubOTAProgress, s.handleOtaProgress)

	// 3. 标准业务事件上报 ➔ 绑定 handleDeviceEvent
	s.subscribeTopic(enum.MQTTSubEvent, s.handleDeviceEvent)

	<-ctx.Done()
	s.client.Disconnect(250)
	s.logger.Info("MQTT client disconnected")
	return nil
}

// subscribeTopic 辅助方法：订阅单个 Topic 并绑定专属 Handler
func (s *Service) subscribeTopic(topic string, handler mqtt.MessageHandler) {
	subToken := s.client.Subscribe(topic, 0, handler)
	if subToken.Wait() && subToken.Error() != nil {
		s.logger.Error("failed to subscribe to topic", zap.String("topic", topic), zap.Error(subToken.Error()))
	} else {
		s.logger.Info("registered dedicated handler for MQTT topic", zap.String("topic", topic))
	}
}

// handleTelemetry 专职处理时序遥测与属性 ➔ 投递至 device.telemetry.v1
func (s *Service) handleTelemetry(_ mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	deviceKey := extractDeviceKey(topic)
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from telemetry topic", zap.String("topic", topic))
		return
	}

	deviceMsg := model.DeviceMessage{
		DeviceKey:   deviceKey,
		Transport:   "mqtt",
		MessageType: "telemetry",
		Payload:     json.RawMessage(msg.Payload()),
		Timestamp:   time.Now().UTC(),
		Headers:     map[string]string{"topic": topic},
	}

	_ = s.producer.SendDeviceMessage(context.Background(), deviceMsg)
}

// handleOtaProgress 专职处理 OTA 固件升级进度 ➔ 投递至 ota.report.v1
func (s *Service) handleOtaProgress(_ mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	deviceKey := extractDeviceKey(topic)
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from OTA topic", zap.String("topic", topic))
		return
	}

	deviceMsg := model.DeviceMessage{
		DeviceKey:   deviceKey,
		Transport:   "mqtt",
		MessageType: "ota_report",
		Payload:     json.RawMessage(msg.Payload()),
		Timestamp:   time.Now().UTC(),
		Headers:     map[string]string{"topic": topic},
	}

	_ = s.producer.SendDeviceMessage(context.Background(), deviceMsg)
}

// handleDeviceEvent 专职处理设备特定告警与事件 ➔ 投递至 device.event.v1
func (s *Service) handleDeviceEvent(_ mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	deviceKey := extractDeviceKey(topic)
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from event topic", zap.String("topic", topic))
		return
	}

	deviceMsg := model.DeviceMessage{
		DeviceKey:   deviceKey,
		Transport:   "mqtt",
		MessageType: "event",
		Payload:     json.RawMessage(msg.Payload()),
		Timestamp:   time.Now().UTC(),
		Headers:     map[string]string{"topic": topic},
	}

	_ = s.producer.SendDeviceMessage(context.Background(), deviceMsg)
}

// HandleDownlinkCommand 处理来自 Kafka 的下行控制指令并推向设备（使用标准枚举模板）
func (s *Service) HandleDownlinkCommand(ctx context.Context, cmd model.DeviceCommand) error {
	topic := cmd.Endpoint
	if topic == "" {
		switch cmd.CommandType {
		case "ota_upgrade", "ota":
			// 使用枚举模板构造 OTA 下发 Topic
			topic = fmt.Sprintf(enum.MQTTTplOTAUpgrade, cmd.DeviceKey)
		default:
			// 使用枚举模板构造属性设置下发 Topic
			topic = fmt.Sprintf(enum.MQTTTplPropertySet, cmd.DeviceKey)
		}
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

// extractDeviceKey 从标准 /sys/{deviceKey}/... 路径中提取 deviceKey
func extractDeviceKey(topic string) string {
	parts := strings.Split(strings.Trim(topic, "/"), "/")
	if len(parts) >= 2 && parts[0] == "sys" {
		return parts[1]
	}
	return ""
}
