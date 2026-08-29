package job

import (
	"context"
	"encoding/json"
	"time"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/transport"
	"aiot-backend/pkg/log"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
)

// DeviceMessageConsumer 消费网关标准化上行消息，维护设备端点在线状态。
type DeviceMessageConsumer struct {
	repo         *repository.ProtocolRepository
	logger       *log.Logger
	client       *kgo.Client
	offlineAfter time.Duration
}

func NewDeviceMessageConsumer(config *viper.Viper, repo *repository.ProtocolRepository, logger *log.Logger) (*DeviceMessageConsumer, error) {
	c := &DeviceMessageConsumer{repo: repo, logger: logger, offlineAfter: 5 * time.Minute}
	if value := config.GetDuration("device_gateway.endpoint_offline_after"); value > 0 {
		c.offlineAfter = value
	}
	brokers := config.GetStringSlice("data.kafka.device.brokers")
	if len(brokers) == 0 {
		return c, nil
	}
	client, err := kgo.NewClient(kgo.SeedBrokers(brokers...), kgo.ConsumerGroup("aiot-backend-device-messages"), kgo.ConsumeTopics(enum.KafkaTopicDeviceMessageV1), kgo.DisableAutoCommit())
	if err != nil {
		return nil, err
	}
	c.client = client
	return c, nil
}

func (c *DeviceMessageConsumer) Start(ctx context.Context) error {
	go c.markOffline(ctx)
	if c.client == nil {
		return nil
	}
	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if len(fetches.Errors()) > 0 {
			continue
		}
		fetches.EachRecord(func(record *kgo.Record) {
			var message transport.DeviceMessage
			if json.Unmarshal(record.Value, &message) != nil || message.DeviceKey == "" {
				_ = c.client.CommitRecords(ctx, record)
				return
			}
			if c.repo != nil {
				if protocol := message.Headers["transport"]; protocol != "" {
					_ = c.repo.MarkEndpointSeenForEvent(ctx, message.DeviceKey, protocol, message.Timestamp)
				}
			}
			_ = c.client.CommitRecords(ctx, record)
		})
	}
	return ctx.Err()
}

func (c *DeviceMessageConsumer) markOffline(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			if c.repo != nil {
				_ = c.repo.MarkStaleEndpointsOffline(ctx, now.Add(-c.offlineAfter))
			}
		}
	}
}

func (c *DeviceMessageConsumer) Stop() {
	if c.client != nil {
		c.client.Close()
	}
}
