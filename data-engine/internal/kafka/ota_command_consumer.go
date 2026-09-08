package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"data-engine/internal/enum"
	"data-engine/internal/model"
	"data-engine/internal/service"
	"data-engine/internal/storage"

	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type OTACommandConsumer struct {
	client     *kgo.Client
	dispatcher *service.OTADispatcher
	store      *storage.OTAStore
	logger     *zap.Logger
}

func NewOTACommandConsumer(config *viper.Viper, logger *zap.Logger, dispatcher *service.OTADispatcher, store *storage.OTAStore) (*OTACommandConsumer, error) {
	brokers := config.GetStringSlice("kafka.brokers")
	if len(brokers) == 0 {
		return &OTACommandConsumer{dispatcher: dispatcher, store: store, logger: logger}, nil
	}
	group := config.GetString("kafka.ota_command_consumer_group")
	if group == "" {
		group = enum.ConsumerGroupOTACommand
	}
	client, err := kgo.NewClient(kgo.SeedBrokers(brokers...), kgo.ConsumerGroup(group), kgo.ConsumeTopics(enum.KafkaTopicOTACommand), kgo.DisableAutoCommit())
	if err != nil {
		return nil, fmt.Errorf("create OTA command consumer: %w", err)
	}
	return &OTACommandConsumer{client: client, dispatcher: dispatcher, store: store, logger: logger}, nil
}

func (c *OTACommandConsumer) Start(ctx context.Context) error {
	if c.client == nil {
		<-ctx.Done()
		return nil
	}
	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				c.logger.Warn("OTA command kafka poll error", zap.Error(err.Err))
			}
			continue
		}
		fetches.EachRecord(func(record *kgo.Record) {
			var command model.OTAUpgradeCommand
			if err := json.Unmarshal(record.Value, &command); err != nil {
				c.logger.Error("invalid OTA command", zap.Error(err))
				c.client.MarkCommitRecords(record)
				return
			}
			var dispatchError string
			if err := c.dispatcher.Dispatch(ctx, command); err != nil {
				dispatchError = err.Error()
			}
			if err := c.store.RecordDispatchResult(ctx, command.BatchID, command.DeviceKey, dispatchError); err != nil {
				c.logger.Error("record OTA dispatch result", zap.Error(err))
				return
			}
			c.client.MarkCommitRecords(record)
		})
		if err := c.client.CommitMarkedOffsets(ctx); err != nil {
			c.logger.Warn("commit OTA command offsets", zap.Error(err))
		}
	}
	return nil
}
func (c *OTACommandConsumer) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
