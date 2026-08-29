package job

import (
	"context"
	"encoding/json"
	"fmt"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/service"
	"aiot-backend/internal/transport"
	"aiot-backend/pkg/log"

	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type otaUpgradeCommand struct {
	BatchID           string `json:"batch_id"`
	ProductKey        string `json:"product_key"`
	DeviceKey         string `json:"device_key"`
	DeviceName        string `json:"device_name"`
	TargetVersion     string `json:"target_version"`
	Module            string `json:"module"`
	URL               string `json:"url"`
	Size              int64  `json:"size"`
	Checksum          string `json:"checksum"`
	TransportProtocol string `json:"transport_protocol"`
	EndpointID        int64  `json:"endpoint_id"`
	Endpoint          string `json:"endpoint"`
}

// OTACommandConsumer 负责消费升级命令并下发到设备 MQTT 主题。
type OTACommandConsumer struct {
	mqtt     service.MQTTServiceInterface
	ota      *service.OTAService
	logger   *log.Logger
	client   *kgo.Client
	registry *transport.Registry
}

func NewOTACommandConsumer(config *viper.Viper, mqtt service.MQTTServiceInterface, ota *service.OTAService, logger *log.Logger, registries ...*transport.Registry) (*OTACommandConsumer, error) {
	c := &OTACommandConsumer{mqtt: mqtt, ota: ota, logger: logger}
	if len(registries) > 0 {
		c.registry = registries[0]
	}
	brokers := config.GetStringSlice("data.kafka.device.brokers")
	if len(brokers) == 0 {
		return c, nil
	}
	client, err := kgo.NewClient(kgo.SeedBrokers(brokers...), kgo.ConsumerGroup("aiot-backend-ota-commands"), kgo.ConsumeTopics(enum.KafkaTopicOTAUpgradeCommandV1), kgo.DisableAutoCommit())
	if err != nil {
		return nil, err
	}
	c.client = client
	return c, nil
}
func (c *OTACommandConsumer) Start(ctx context.Context) error {
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
			var command otaUpgradeCommand
			if err := json.Unmarshal(record.Value, &command); err != nil {
				c.logger.Error("invalid OTA command message", zap.Error(err))
				_ = c.client.CommitRecords(ctx, record)
				return
			}
			// 基础字段与依赖校验：若消息本身损坏或关键属性缺失，记录日志后提交位点以跳过坏消息（防 Poison Pill 阻塞分区）
			if command.ProductKey == "" || command.DeviceKey == "" || command.BatchID == "" {
				c.logger.Error("OTA command missing required fields, skipping record",
					zap.String("batch_id", command.BatchID),
					zap.String("product_key", command.ProductKey),
					zap.String("device_key", command.DeviceKey),
				)
				_ = c.client.CommitRecords(ctx, record)
				return
			}
			if c.mqtt == nil || c.ota == nil {
				c.logger.Error("OTA command consumer dependency unavailable (mqtt or ota service nil)")
				return
			}

			// 先原子领取任务再发布 MQTT，避免 Kafka 重放或并发消费导致重复下发。
			claimed, err := c.ota.ClaimBatchDeviceForMQTT(ctx, command.BatchID, command.DeviceKey)
			if err != nil {
				c.logger.Error("failed to claim OTA command", zap.Error(err))
				return
			}
			if !claimed {
				// 已领取或已终态的任务直接确认位点，安全跳过重复消息。
				_ = c.client.CommitRecords(ctx, record)
				return
			}

			// 组装统一下行命令；具体协议由 Transport Registry 选择。
			payload, err := json.Marshal(command)
			if err == nil {
				deviceName := command.DeviceName
				if deviceName == "" {
					deviceName = command.DeviceKey
				}
				topic := fmt.Sprintf(enum.MQTTTopicOTADeviceUpgrade, command.ProductKey, deviceName)
				protocol := command.TransportProtocol
				if protocol == "" {
					protocol = string(enum.TransportMQTT)
				}
				if c.registry != nil {
					adapter, ok := c.registry.Get(enum.TransportProtocol(protocol))
					if !ok {
						err = fmt.Errorf("no adapter registered for transport protocol %q", protocol)
					} else {
						err = adapter.Send(ctx, transport.Command{DeviceKey: command.DeviceKey, EndpointID: command.EndpointID, Type: "ota", Payload: payload, Headers: map[string]string{"topic": topic, "endpoint": command.Endpoint}})
					}
				} else {
					// 使用 QoS 1 确保旧部署仍能通过 MQTT Broker 下发。
					err = c.mqtt.Publish(ctx, topic, 1, payload)
				}
			}
			if err != nil {
				c.logger.Error("failed to publish OTA command to MQTT", zap.Error(err))
				if resetErr := c.ota.ResetMQTTDispatch(ctx, command.BatchID, command.DeviceKey, err.Error()); resetErr != nil {
					c.logger.Error("failed to reset OTA command after MQTT publish failure", zap.Error(resetErr))
				}
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
