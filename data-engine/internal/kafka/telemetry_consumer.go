package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"data-engine/internal/engine"
	"data-engine/internal/enum"
	"data-engine/internal/model"

	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

// TelemetryConsumer 封装 Kafka 消费者客户端，专职订阅并计算 device.telemetry.v1 遥测主题。
type TelemetryConsumer struct {
	client    *kgo.Client
	processor *engine.Processor
	logger    *zap.Logger
}

// NewTelemetryConsumer 初始化遥测流消费者。
func NewTelemetryConsumer(config *viper.Viper, logger *zap.Logger, processor *engine.Processor) (*TelemetryConsumer, error) {
	brokers := config.GetStringSlice("kafka.brokers")
	if len(brokers) == 0 {
		return &TelemetryConsumer{processor: processor, logger: logger}, nil
	}

	topic := enum.KafkaTopicDeviceTelemetry

	group := config.GetString("kafka.telemetry_consumer_group")
	if group == "" {
		group = enum.ConsumerGroupTelemetry
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.DisableAutoCommit(),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry kafka consumer: %w", err)
	}

	return &TelemetryConsumer{
		client:    client,
		processor: processor,
		logger:    logger,
	}, nil
}

// Start 启动遥测流拉取与计算循环。
func (c *TelemetryConsumer) Start(ctx context.Context) error {
	if c.client == nil {
		c.logger.Info("telemetry kafka consumer running in mock/disabled mode")
		<-ctx.Done()
		return nil
	}

	c.logger.Info("Data Engine Telemetry Consumer started polling Kafka records...")

	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				c.logger.Warn("telemetry kafka poll error", zap.Error(err.Err), zap.String("topic", err.Topic))
			}
			continue
		}

		fetches.EachRecord(func(record *kgo.Record) {
			var msg model.DeviceMessage
			if err := json.Unmarshal(record.Value, &msg); err != nil {
				c.logger.Error("failed to unmarshal telemetry message (poison pill skipped)", zap.Error(err))
				c.client.MarkCommitRecords(record)
				return
			}

			if err := c.processor.ProcessMessage(ctx, msg); err != nil {
				c.logger.Error("failed to process telemetry message", zap.String("device_key", msg.DeviceKey), zap.Error(err))
			}

			c.client.MarkCommitRecords(record)
		})

		// 触发持久化提交到 Kafka Broker
		if err := c.client.CommitMarkedOffsets(ctx); err != nil {
			c.logger.Warn("failed to commit telemetry kafka offsets", zap.Error(err))
		}
	}
	return nil
}

// Close 安全关闭消费者连接
func (c *TelemetryConsumer) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
