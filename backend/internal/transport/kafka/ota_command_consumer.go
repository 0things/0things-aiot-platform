package kafkatransport

import (
	"context"
	"encoding/json"
	"fmt"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/transport"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
)

type otaCommand struct {
	BatchID           string          `json:"batch_id"`
	ProductKey        string          `json:"product_key"`
	DeviceKey         string          `json:"device_key"`
	DeviceName        string          `json:"device_name"`
	TransportProtocol string          `json:"transport_protocol"`
	EndpointID        int64           `json:"endpoint_id"`
	Endpoint          string          `json:"endpoint"`
	Payload           json.RawMessage `json:"-"`
}

// OTACommandConsumer 在设备网关内消费 OTA 指令，并交给协议适配器下发。
// 与管理服务使用同一 consumer group，确保同一命令不会被两个进程重复处理。
type OTACommandConsumer struct {
	client   *kgo.Client
	registry *transport.Registry
}

func NewOTACommandConsumer(config *viper.Viper, registry *transport.Registry) (*OTACommandConsumer, error) {
	brokers := config.GetStringSlice("data.kafka.device.brokers")
	c := &OTACommandConsumer{registry: registry}
	if len(brokers) == 0 {
		return c, nil
	}
	group := config.GetString("device_gateway.kafka_consumer_group")
	if group == "" {
		group = "aiot-backend-ota-commands"
	}
	client, err := kgo.NewClient(kgo.SeedBrokers(brokers...), kgo.ConsumerGroup(group), kgo.ConsumeTopics(enum.KafkaTopicOTAUpgradeCommandV1), kgo.DisableAutoCommit())
	if err != nil {
		return nil, err
	}
	c.client = client
	return c, nil
}

func (c *OTACommandConsumer) Start(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			continue
		}
		fetches.EachRecord(func(record *kgo.Record) {
			var command otaCommand
			if err := json.Unmarshal(record.Value, &command); err != nil {
				_ = c.client.CommitRecords(ctx, record)
				return
			}
			if command.DeviceKey == "" || command.ProductKey == "" || c.registry == nil {
				_ = c.client.CommitRecords(ctx, record)
				return
			}
			protocol := command.TransportProtocol
			if protocol == "" {
				protocol = string(enum.TransportMQTT)
			}
			adapter, ok := c.registry.Get(enum.TransportProtocol(protocol))
			if !ok {
				err := fmt.Errorf("no adapter registered for transport protocol %q", protocol)
				c.reportFailure(ctx, command, err)
				_ = c.client.CommitRecords(ctx, record)
				return
			}
			deviceName := command.DeviceName
			if deviceName == "" {
				deviceName = command.DeviceKey
			}
			topic := fmt.Sprintf(enum.MQTTTopicOTADeviceUpgrade, command.ProductKey, deviceName)
			if err := adapter.Send(ctx, transport.Command{DeviceKey: command.DeviceKey, EndpointID: command.EndpointID, Type: "ota", Payload: record.Value, Headers: map[string]string{"topic": topic, "endpoint": command.Endpoint}}); err != nil {
				c.reportFailure(ctx, command, err)
				_ = c.client.CommitRecords(ctx, record)
				return
			}
			_ = c.client.CommitRecords(ctx, record)
		})
	}
	return ctx.Err()
}

func (c *OTACommandConsumer) reportFailure(ctx context.Context, command otaCommand, err error) {
	payload, marshalErr := json.Marshal(map[string]any{
		"batch_id":    command.BatchID,
		"device_key":  command.DeviceKey,
		"status":      "failed",
		"progress":    0,
		"description": err.Error(),
	})
	if marshalErr != nil {
		return
	}
	result := c.client.ProduceSync(ctx, &kgo.Record{Topic: enum.KafkaTopicOTAUpgradeReportV1, Key: []byte(command.DeviceKey), Value: payload})
	if result.FirstErr() != nil {
		// 下发失败时无法再更新状态，只保留原始错误让消费者重试。
		return
	}
}

func (c *OTACommandConsumer) Stop() {
	if c.client != nil {
		c.client.Close()
	}
}
