package job

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/service"
	"aiot-backend/pkg/log"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type otaUpgradeCommand struct {
	BatchID       string `json:"batch_id"`
	ProductKey    string `json:"product_key"`
	DeviceKey     string `json:"device_key"`
	DeviceName    string `json:"device_name"`
	TargetVersion string `json:"target_version"`
	Module        string `json:"module"`
	URL           string `json:"url"`
	Size          int64  `json:"size"`
	Checksum      string `json:"checksum"`
}

type OTACommandConsumer struct {
	mqtt   service.MQTTServiceInterface
	ota    *service.OTAService
	kafka  service.KafkaServiceInterface
	logger *log.Logger
	client *kgo.Client
}

func NewOTACommandConsumer(config *viper.Viper, mqtt service.MQTTServiceInterface, ota *service.OTAService, kafka service.KafkaServiceInterface, logger *log.Logger) (*OTACommandConsumer, error) {
	brokers := config.GetStringSlice("data.kafka.device.brokers")
	consumer := &OTACommandConsumer{mqtt: mqtt, ota: ota, kafka: kafka, logger: logger}
	if len(brokers) == 0 {
		return consumer, nil
	}
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup("aiot-backend-ota-commands"),
		kgo.ConsumeTopics(enum.KafkaTopicOTAUpgradeCommandV1, enum.KafkaTopicOTAUpgradeReportV1),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		return nil, err
	}
	consumer.client = client
	return consumer, nil
}

func (c *OTACommandConsumer) Start(ctx context.Context) error {
	// MQTT 只负责设备侧传输，设备回报统一转发到 Kafka，再由本消费者落库。
	if c.mqtt != nil && c.ota != nil {
		handler := func(_ mqtt.Client, msg mqtt.Message) {
			var report struct {
				Params struct {
					BatchID   string          `json:"batch_id"`
					DeviceKey string          `json:"device_key"`
					Status    string          `json:"status"`
					Step      json.RawMessage `json:"step"`
					Version   string          `json:"version"`
					Desc      string          `json:"desc"`
				} `json:"params"`
			}
			if err := json.Unmarshal(msg.Payload(), &report); err != nil {
				c.logger.Error("invalid OTA MQTT report", zap.Error(err))
				return
			}
			if report.Params.BatchID == "" || report.Params.DeviceKey == "" {
				c.logger.Warn("OTA MQTT report missing batch_id or device_key")
				return
			}
			status := report.Params.Status
			if status == "" {
				status = enum.OTAStatusInProgress
			}
			var step int32
			if err := json.Unmarshal(report.Params.Step, &step); err != nil {
				var stepText string
				if json.Unmarshal(report.Params.Step, &stepText) == nil {
					stepValue, _ := strconv.ParseInt(stepText, 10, 32)
					step = int32(stepValue)
				}
			}
			if c.kafka == nil {
				c.logger.Warn("OTA report Kafka producer unavailable")
				return
			}
			if err := c.kafka.ProduceJSON(ctx, enum.KafkaTopicOTAUpgradeReportV1, report.Params.BatchID+":"+report.Params.DeviceKey, map[string]any{
				"batch_id": report.Params.BatchID, "device_key": report.Params.DeviceKey,
				"status": status, "version": report.Params.Version, "progress": step, "desc": report.Params.Desc,
			}); err != nil {
				c.logger.Error("failed to publish OTA MQTT report", zap.Error(err), zap.String("desc", report.Params.Desc))
			}
		}
		for _, topic := range []string{enum.MQTTTopicOTADeviceProgress, enum.MQTTTopicOTADeviceInform} {
			if err := c.mqtt.Subscribe(ctx, topic, 1, handler); err != nil {
				c.logger.Warn("failed to subscribe OTA MQTT topic", zap.String("topic", topic), zap.Error(err))
			}
		}
	}
	if c.client == nil {
		c.logger.Info("OTA command Kafka consumer disabled: no brokers configured")
		return nil
	}
	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				c.logger.Error("OTA command Kafka poll failed", zap.Error(err.Err))
			}
			continue
		}
		fetches.EachRecord(func(record *kgo.Record) {
			if record.Topic == enum.KafkaTopicOTAUpgradeReportV1 {
				var report struct {
					BatchID   string `json:"batch_id"`
					DeviceKey string `json:"device_key"`
					Status    string `json:"status"`
					Version   string `json:"version"`
					Progress  int32  `json:"progress"`
				}
				if err := json.Unmarshal(record.Value, &report); err != nil {
					c.logger.Error("invalid OTA report message", zap.Error(err))
					_ = c.client.CommitRecords(ctx, record)
					return
				}
				if err := c.ota.ReportBatchDevice(ctx, report.BatchID, report.DeviceKey, report.Status, report.Version, report.Progress); err != nil {
					c.logger.Error("failed to process OTA report", zap.Error(err))
					return
				}
				if err := c.client.CommitRecords(ctx, record); err != nil {
					c.logger.Error("failed to commit OTA report", zap.Error(err))
				}
				return
			}
			var command otaUpgradeCommand
			if err := json.Unmarshal(record.Value, &command); err != nil {
				c.logger.Error("invalid OTA command message", zap.Error(err))
				_ = c.client.CommitRecords(ctx, record)
				return
			}
			if command.ProductKey == "" || command.DeviceKey == "" || command.BatchID == "" {
				c.logger.Error("OTA command missing required fields")
				_ = c.client.CommitRecords(ctx, record)
				return
			}
			shouldDispatch, err := c.ota.ShouldDispatchBatchDevice(ctx, command.BatchID, command.DeviceKey)
			if err != nil {
				c.logger.Error("failed to check OTA command state", zap.Error(err))
				return
			}
			if !shouldDispatch {
				if err := c.client.CommitRecords(ctx, record); err != nil {
					c.logger.Error("failed to commit duplicate OTA command", zap.Error(err))
				}
				return
			}
			payload, err := json.Marshal(command)
			if err == nil && c.mqtt != nil {
				deviceName := command.DeviceName
				if deviceName == "" {
					deviceName = command.DeviceKey
				}
				topic := fmt.Sprintf(enum.MQTTTopicOTADeviceUpgrade, command.ProductKey, deviceName)
				err = c.mqtt.Publish(ctx, topic, 1, payload)
			}
			if err != nil {
				c.logger.Error("failed to publish OTA command to MQTT", zap.Error(err))
				return
			}
			if err := c.client.CommitRecords(ctx, record); err != nil {
				c.logger.Error("failed to commit OTA command", zap.Error(err))
			}
		})
	}
	return ctx.Err()
}

func (c *OTACommandConsumer) Stop() {
	if c.client != nil {
		c.client.Close()
	}
}
