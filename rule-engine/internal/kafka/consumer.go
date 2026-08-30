package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"rule-engine/internal/engine"
	"rule-engine/internal/enum"
	"rule-engine/internal/model"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

// Consumer 封装 Kafka 消费者客户端，作为独立的计算 Worker 订阅 device.telemetry.v1 遥测主题。
type Consumer struct {
	client    *kgo.Client
	processor *engine.Processor
	logger    *zap.Logger
}

// NewConsumer 初始化规则引擎 Kafka 消费组。
// 约定优于配置：消费组优先读取配置文件中的 kafka.consumer_group，若未配置则回退到代码枚举常量。
func NewConsumer(config *viper.Viper, logger *zap.Logger, processor *engine.Processor) (*Consumer, error) {
	brokers := config.GetStringSlice("kafka.brokers")
	if len(brokers) == 0 {
		return &Consumer{processor: processor, logger: logger}, nil
	}

	topic := enum.KafkaTopicDeviceTelemetry

	group := config.GetString("kafka.consumer_group")
	if group == "" {
		group = enum.ConsumerGroupRuleEngine
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.DisableAutoCommit(), // 关闭自动提交，处理成功后手动确认
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer client: %w", err)
	}

	return &Consumer{
		client:    client,
		processor: processor,
		logger:    logger,
	}, nil
}

// Start 启动消息拉取主循环，将拉取到的 Record 交由 Processor 处理，成功后提交 Offset。
func (c *Consumer) Start(ctx context.Context) error {
	if c.client == nil {
		c.logger.Info("kafka consumer running in standalone/mock mode")
		<-ctx.Done()
		return nil
	}

	c.logger.Info("Rule Engine started polling Kafka records...")

	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			continue
		}

		fetches.EachRecord(func(record *kgo.Record) {
			var msg model.DeviceMessage
			if err := json.Unmarshal(record.Value, &msg); err != nil {
				c.logger.Error("failed to unmarshal device message", zap.Error(err))
				return
			}

			// 调用规则计算处理器
			if err := c.processor.ProcessMessage(ctx, msg); err != nil {
				c.logger.Error("failed to process device message", zap.String("device_key", msg.DeviceKey), zap.Error(err))
				return
			}

			// 处理完毕后确认提交位移
			c.client.MarkCommitRecords(record)
		})
	}
	return nil
}

// Close 安全关闭 Kafka 消费者
func (c *Consumer) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
