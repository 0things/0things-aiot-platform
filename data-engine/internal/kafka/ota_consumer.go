package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"data-engine/internal/enum"
	"data-engine/internal/model"
	"data-engine/internal/service"

	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

// OTAConsumer 封装专用于消费 OTA 进度的 Kafka Consumer。
type OTAConsumer struct {
	client    *kgo.Client
	processor *service.OTAProcessor
	logger    *zap.Logger
}

// NewOTAConsumer 初始化 OTA 进度消费者。
func NewOTAConsumer(config *viper.Viper, logger *zap.Logger, processor *service.OTAProcessor) (*OTAConsumer, error) {
	brokers := config.GetStringSlice("kafka.brokers")
	if len(brokers) == 0 {
		return &OTAConsumer{processor: processor, logger: logger}, nil
	}

	// 专职只监听 OTA 进度主题
	topic := enum.KafkaTopicOTAReport

	group := config.GetString("kafka.ota_consumer_group")
	if group == "" {
		group = enum.ConsumerGroupOTA
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

	return &OTAConsumer{
		client:    client,
		processor: processor,
		logger:    logger,
	}, nil
}

// Start 启动 OTA 进度消费循环。
func (c *OTAConsumer) Start(ctx context.Context) error {
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
			report, ok := c.extractOTAReport(record.Value)
			if !ok {
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

// extractOTAReport 兼容标准 OTADeviceUpgradeReportEvent 与 通用 DeviceMessage(MessageType="ota_report")。
func (c *OTAConsumer) extractOTAReport(data []byte) (model.OTADeviceUpgradeReportEvent, bool) {
	var report model.OTADeviceUpgradeReportEvent
	if err := json.Unmarshal(data, &report); err == nil && report.DeviceKey != "" && (report.Progress != 0 || report.Status != "") {
		return report, true
	}

	var deviceMsg model.DeviceMessage
	if err := json.Unmarshal(data, &deviceMsg); err == nil {
		if deviceMsg.MessageType == "ota_report" || deviceMsg.MessageType == "ota_progress" {
			var payloadMap map[string]interface{}
			if err := json.Unmarshal(deviceMsg.Payload, &payloadMap); err == nil {
				rep := model.OTADeviceUpgradeReportEvent{
					DeviceKey: deviceMsg.DeviceKey,
					Timestamp: deviceMsg.Timestamp,
				}
				if step, ok := payloadMap["step"]; ok {
					switch v := step.(type) {
					case float64:
						rep.Progress = int(v)
					case string:
						rep.Progress, _ = strconv.Atoi(v)
					}
				}
				if desc, ok := payloadMap["desc"].(string); ok {
					rep.Desc = desc
				}
				if status, ok := payloadMap["status"].(string); ok {
					rep.Status = status
				}
				if batchID, ok := payloadMap["batch_id"].(string); ok {
					rep.BatchID = batchID
				}
				return rep, true
			}
		}
	}

	return model.OTADeviceUpgradeReportEvent{}, false
}

// Close 安全关闭 Kafka 消费者连接
func (c *OTAConsumer) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
