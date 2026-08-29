package job

import (
	"context"
	"encoding/json"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/service"
	"aiot-backend/pkg/log"

	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

// OTAReportConsumer 负责消费 Kafka 回报并更新 OTA 批次和设备状态。
type OTAReportConsumer struct {
	ota    *service.OTAService
	logger *log.Logger
	client *kgo.Client
}

func NewOTAReportConsumer(config *viper.Viper, ota *service.OTAService, logger *log.Logger) (*OTAReportConsumer, error) {
	c := &OTAReportConsumer{ota: ota, logger: logger}
	brokers := config.GetStringSlice("data.kafka.device.brokers")
	if len(brokers) == 0 {
		return c, nil
	}
	client, err := kgo.NewClient(kgo.SeedBrokers(brokers...), kgo.ConsumerGroup("aiot-backend-ota-reports"), kgo.ConsumeTopics(enum.KafkaTopicOTAUpgradeReportV1), kgo.DisableAutoCommit())
	if err != nil {
		return nil, err
	}
	c.client = client
	return c, nil
}
func (c *OTAReportConsumer) Start(ctx context.Context) error {
	if c.client == nil {
		c.logger.Info("OTA report Kafka consumer disabled: no brokers configured")
		return nil
	}
	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				c.logger.Error("OTA report Kafka poll failed", zap.Error(err.Err))
			}
			continue
		}
		fetches.EachRecord(func(record *kgo.Record) {
			var report OTADeviceUpgradeReportEvent
			if err := json.Unmarshal(record.Value, &report); err != nil {
				c.logger.Error("invalid OTA report message, skipping record", zap.Error(err))
				_ = c.client.CommitRecords(ctx, record)
				return
			}

			// 畸形消息校验：缺少批次或设备标识直接跳过并确认位点
			if report.BatchID == "" || report.DeviceKey == "" {
				c.logger.Warn("OTA report missing required batch_id or device_key, skipping record")
				_ = c.client.CommitRecords(ctx, record)
				return
			}

			if c.ota == nil {
				c.logger.Error("OTA service dependency unavailable")
				return
			}

			// 调用业务层处理上报结果，同时持久化可能的错误原因（report.Desc）
			if err := c.ota.ReportBatchDevice(ctx, report.BatchID, report.DeviceKey, report.Status, report.Version, report.Progress, report.Desc); err != nil {
				c.logger.Error("failed to process OTA report",
					zap.Error(err),
					zap.String("batch_id", report.BatchID),
					zap.String("device_key", report.DeviceKey),
				)
				return
			}
			_ = c.client.CommitRecords(ctx, record)
		})
	}
	return ctx.Err()
}
func (c *OTAReportConsumer) Stop() {
	if c.client != nil {
		c.client.Close()
	}
}
