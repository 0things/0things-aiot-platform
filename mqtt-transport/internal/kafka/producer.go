package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"mqtt-transport/internal/enum"
	"mqtt-transport/internal/model"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

// Producer 封装 Franz-Go Kafka 生产者，根据消息类型（遥测/OTA/事件）将报文动态路由至专属 Kafka Topic。
type Producer struct {
	client *kgo.Client
	logger *zap.Logger
}

// NewProducer 初始化 Kafka 分流生产者。
func NewProducer(config *viper.Viper, logger *zap.Logger) (*Producer, func(), error) {
	brokers := config.GetStringSlice("kafka.brokers")
	if len(brokers) == 0 {
		logger.Warn("kafka brokers not configured, running in standalone/mock mode")
		return &Producer{logger: logger}, func() {}, nil
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create kafka producer client: %w", err)
	}

	cleanup := func() {
		client.Close()
	}

	return &Producer{
		client: client,
		logger: logger,
	}, cleanup, nil
}

// SendDeviceMessage 根据消息的 MessageType 动态选择投递的目标 Kafka Topic。
func (p *Producer) SendDeviceMessage(ctx context.Context, msg model.DeviceMessage) error {
	var targetTopic string
	switch msg.MessageType {
	case "ota_report", "ota_progress":
		targetTopic = enum.KafkaTopicOTAReport
	case "event":
		targetTopic = enum.KafkaTopicDeviceEvent
	default:
		targetTopic = enum.KafkaTopicDeviceTelemetry
	}

	if p.client == nil {
		p.logger.Debug("mock kafka produce",
			zap.String("target_topic", targetTopic),
			zap.String("device_key", msg.DeviceKey),
			zap.String("type", msg.MessageType),
		)
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal device message: %w", err)
	}

	record := &kgo.Record{
		Topic: targetTopic,
		Key:   []byte(msg.DeviceKey), // 按 DeviceKey 分区保证单设备严格时序
		Value: data,
	}

	return p.client.ProduceSync(ctx, record).FirstErr()
}
