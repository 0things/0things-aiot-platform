package job

import (
	"context"
	"encoding/json"
	"errors"

	"aiot-backend/internal/repository"
	"aiot-backend/internal/service"
	"aiot-backend/pkg/log"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

const deviceEventTopic = "iotEvent"

type deviceEventMessage struct {
	Type       string         `json:"type"`
	Timestamp  int64          `json:"timestamp"`
	ProductKey string         `json:"product_key"`
	DeviceKey  string         `json:"device_key"`
	Data       map[string]any `json:"data"`
}

type DeviceEventConsumer struct {
	service *service.DeviceEventService
	logger  *log.Logger
	client  *kgo.Client
}

func NewDeviceEventConsumer(config *viper.Viper, service *service.DeviceEventService, logger *log.Logger) (*DeviceEventConsumer, error) {
	brokers := config.GetStringSlice("data.kafka.device.brokers")
	if len(brokers) == 0 { return &DeviceEventConsumer{service: service, logger: logger}, nil }
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup("aiot-backend-device-events"),
		kgo.ConsumeTopics(deviceEventTopic),
		kgo.DisableAutoCommit(),
	)
	if err != nil { return nil, err }
	return &DeviceEventConsumer{service: service, logger: logger, client: client}, nil
}

func (c *DeviceEventConsumer) Start(ctx context.Context) {
	if c.client == nil { c.logger.Info("device event Kafka consumer disabled: no brokers configured"); return }
	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				c.logger.Error("device event Kafka poll failed", zap.Error(err.Err))
			}
			continue
		}
		fetches.EachRecord(func(record *kgo.Record) {
			var event deviceEventMessage
			if err := json.Unmarshal(record.Value, &event); err != nil {
				c.logger.Error("invalid device event message", zap.Error(err))
				_ = c.client.CommitRecords(ctx, record)
				return
			}
			if err := c.service.Record(ctx, event.ProductKey, event.DeviceKey, event.Type, event.Timestamp, event.Data); err != nil {
				c.logger.Error("failed to store device event", zap.Error(err))
				if errors.Is(err, service.ErrInvalidDeviceEvent) || errors.Is(err, repository.ErrNotFound) {
					_ = c.client.CommitRecords(ctx, record)
				}
				return
			}
			if err := c.client.CommitRecords(ctx, record); err != nil {
				c.logger.Error("failed to commit device event", zap.Error(err))
			}
		})
	}
}

func (c *DeviceEventConsumer) Stop() { if c.client != nil { c.client.Close() } }
