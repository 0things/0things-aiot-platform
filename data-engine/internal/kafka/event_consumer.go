package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"data-engine/internal/enum"
	"data-engine/internal/model"
	"data-engine/internal/service"

	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

// EventConsumer 封装专用于消费设备事件 (device.event.v1) 的 Kafka Consumer。
type EventConsumer struct {
	client    *kgo.Client
	processor *service.EventProcessor
	logger    *zap.Logger
}

// NewEventConsumer 初始化事件流消费者。
func NewEventConsumer(config *viper.Viper, logger *zap.Logger, processor *service.EventProcessor) (*EventConsumer, error) {
	brokers := config.GetStringSlice("kafka.brokers")
	if len(brokers) == 0 {
		return &EventConsumer{processor: processor, logger: logger}, nil
	}

	topic := enum.KafkaTopicDeviceEvent

	group := config.GetString("kafka.event_consumer_group")
	if group == "" {
		group = enum.ConsumerGroupEvent
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.DisableAutoCommit(),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create event kafka consumer client: %w", err)
	}

	return &EventConsumer{
		client:    client,
		processor: processor,
		logger:    logger,
	}, nil
}

// Start 启动设备事件流拉取与处理循环。
func (c *EventConsumer) Start(ctx context.Context) error {
	if c.client == nil {
		c.logger.Info("Event consumer running in standalone/mock mode")
		<-ctx.Done()
		return nil
	}

	c.logger.Info("Data Engine Event Consumer started polling Kafka records...")

	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			continue
		}

		fetches.EachRecord(func(record *kgo.Record) {
			var msg model.DeviceMessage
			if err := json.Unmarshal(record.Value, &msg); err != nil {
				c.logger.Error("failed to unmarshal device event message", zap.Error(err))
				c.client.MarkCommitRecords(record)
				return
			}

			if err := c.processor.HandleEvent(ctx, msg); err != nil {
				c.logger.Error("failed to process device event", zap.String("device_key", msg.DeviceKey), zap.Error(err))
				return
			}

			c.client.MarkCommitRecords(record)
		})
	}
	return nil
}

// Close 安全关闭 Kafka 消费者连接
func (c *EventConsumer) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
