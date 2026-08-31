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

// DownlinkHandler 定义下行控制指令的处理函数类型。
type DownlinkHandler func(ctx context.Context, cmd model.DeviceCommand) error

// Consumer 负责订阅 Kafka 下行命令主题（device.command.v1），将云端下发的控制指令交由传输协议网关推给终端。
type Consumer struct {
	client  *kgo.Client
	handler DownlinkHandler
	logger  *zap.Logger
}

// NewConsumer 初始化 Kafka 下行消费者。
// 约定优于配置：消费组优先读取配置文件中的 kafka.consumer_group，未配置时使用代码枚举默认值。
func NewConsumer(config *viper.Viper, logger *zap.Logger, handler DownlinkHandler) (*Consumer, error) {
	brokers := config.GetStringSlice("kafka.brokers")
	if len(brokers) == 0 {
		return &Consumer{handler: handler, logger: logger}, nil
	}

	topic := enum.KafkaTopicDeviceCommand

	group := config.GetString("kafka.consumer_group")
	if group == "" {
		group = enum.ConsumerGroupMqttDownlink
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.DisableAutoCommit(), // 关闭自动提交，采用手动确认机制
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer client: %w", err)
	}

	return &Consumer{
		client:  client,
		handler: handler,
		logger:  logger,
	}, nil
}

// Start 启动下行拉取长轮询循环。
func (c *Consumer) Start(ctx context.Context) error {
	if c.client == nil {
		c.logger.Info("kafka consumer running in mock/disabled mode")
		<-ctx.Done()
		return nil
	}

	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				c.logger.Warn("kafka poll fetch error", zap.Error(err.Err), zap.String("topic", err.Topic))
			}
			continue
		}

		fetches.EachRecord(func(record *kgo.Record) {
			var cmd model.DeviceCommand
			if err := json.Unmarshal(record.Value, &cmd); err != nil {
				c.logger.Error("failed to unmarshal device command (poison pill skipped)", zap.Error(err))
				c.client.MarkCommitRecords(record)
				return
			}

			// 协议过滤
			if cmd.Transport != "" && cmd.Transport != "mqtt" {
				c.client.MarkCommitRecords(record)
				return
			}

			// 业务分发
			if c.handler != nil {
				if err := c.handler(ctx, cmd); err != nil {
					c.logger.Error("failed to handle downlink command", zap.String("device_key", cmd.DeviceKey), zap.Error(err))
					// 即使单条下发失败，也标记并提交，避免由于终端离线导致整批阻塞
				}
			}

			// 手动标记提交消费位移
			c.client.MarkCommitRecords(record)
		})

		// 触发持久化提交到 Kafka Broker
		if err := c.client.CommitMarkedOffsets(ctx); err != nil {
			c.logger.Warn("failed to commit kafka offsets", zap.Error(err))
		}
	}
	return nil
}

// Close 安全关闭 Kafka 消费者连接
func (c *Consumer) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
