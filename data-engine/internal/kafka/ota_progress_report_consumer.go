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

// OTAProgressReportConsumer consumes device progress and inform reports from
// ota.upgrade.report.v1.
type OTAProgressReportConsumer struct {
	client    *kgo.Client
	processor *service.OTAProcessor
	logger    *zap.Logger
}

func NewOTAProgressReportConsumer(config *viper.Viper, logger *zap.Logger, processor *service.OTAProcessor) (*OTAProgressReportConsumer, error) {
	brokers := config.GetStringSlice("kafka.brokers")
	if len(brokers) == 0 {
		return &OTAProgressReportConsumer{processor: processor, logger: logger}, nil
	}

	topic := enum.KafkaTopicOTAReport

	group := config.GetString("kafka.ota_report_consumer_group")
	if group == "" {
		group = enum.ConsumerGroupOTAReport
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.DisableAutoCommit(),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create ota kafka consumer client: %w", err)
	}

	return &OTAProgressReportConsumer{
		client:    client,
		processor: processor,
		logger:    logger,
	}, nil
}

// Start 启动 OTA 进度消费循环。
func (c *OTAProgressReportConsumer) Start(ctx context.Context) error {
	if c.client == nil {
		c.logger.Info("OTA consumer running in standalone/mock mode")
		<-ctx.Done()
		return nil
	}

	c.logger.Info("Data Engine OTA Consumer started polling Kafka records...")

	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				c.logger.Warn("OTA kafka poll error", zap.Error(err.Err), zap.String("topic", err.Topic))
			}
			continue
		}

		fetches.EachRecord(func(record *kgo.Record) {
			var report model.OTAUpgradeReport
			if err := json.Unmarshal(record.Value, &report); err != nil {
				c.logger.Error("invalid OTA report", zap.Error(err))
				c.client.MarkCommitRecords(record)
				return
			}

			if err := c.processor.HandleOTAReport(ctx, report); err != nil {
				c.logger.Error("failed to process OTA report", zap.String("device_key", report.DeviceKey), zap.Error(err))
			}

			c.client.MarkCommitRecords(record)
		})

		// 触发持久化提交到 Kafka Broker
		if err := c.client.CommitMarkedOffsets(ctx); err != nil {
			c.logger.Warn("failed to commit OTA kafka offsets", zap.Error(err))
		}
	}
	return nil
}

// Close 安全关闭 Kafka 消费者连接
func (c *OTAProgressReportConsumer) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
