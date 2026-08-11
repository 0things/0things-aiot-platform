package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"telemetry-service/internal/biz"
	"telemetry-service/internal/conf"
	"telemetry-service/internal/data/kafka"

	"github.com/go-kratos/kratos/v2/log"
)

const defaultEventConcurrency = 10

// EventConsumer handles consuming events from Kafka.
type EventConsumer struct {
	uc       *biz.EventUsecase
	consumer *kafka.TopicConsumer
	logger   *log.Helper
}

// NewEventConsumer creates a new event consumer.
func NewEventConsumer(uc *biz.EventUsecase, c *conf.Data, logger log.Logger) (*EventConsumer, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "consumer/event"))

	ec := &EventConsumer{
		uc:     uc,
		logger: helper,
	}

	// Create topic consumer
	consumer, cleanup, err := kafka.NewTopicConsumer(c, kafka.TopicConsumerConfig{
		Topic:       kafka.TopicEvent,
		GroupID:     kafka.GroupEvent,
		Handler:     ec.handleEvent,
		Concurrency: defaultEventConcurrency,
	}, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create event consumer: %w", err)
	}

	ec.consumer = consumer

	if consumer != nil {
		helper.Infof("EventConsumer initialized for topic: %s, group: %s", kafka.TopicEvent, kafka.GroupEvent)
	} else {
		helper.Info("Kafka not configured, EventConsumer running in no-op mode")
	}

	return ec, cleanup, nil
}

// handleEvent processes events from Kafka.
func (ec *EventConsumer) handleEvent(ctx context.Context, topic string, key, value []byte) error {
	ec.logger.Debugf("Received event from topic %s with key: %s", topic, string(key))

	// Step 1: Deserialize JSON
	var event biz.Event
	if err := json.Unmarshal(value, &event); err != nil {
		ec.logger.Errorf("Failed to unmarshal event: %v", err)
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Step 2: Basic validation
	if event.Type == "" {
		ec.logger.Error("Event type is required")
		return fmt.Errorf("type is required")
	}
	if event.Timestamp == 0 {
		ec.logger.Error("Timestamp is required")
		return fmt.Errorf("timestamp is required")
	}
	if event.ProductKey == "" {
		ec.logger.Error("Product key is required")
		return fmt.Errorf("product_key is required")
	}
	if event.DeviceKey == "" {
		ec.logger.Error("Device key is required")
		return fmt.Errorf("device_key is required")
	}

	// Step 3: Pass to business logic layer
	if err := ec.uc.ProcessEvent(ctx, &event); err != nil {
		ec.logger.Errorf("Failed to process event: %v", err)
		return err
	}

	return nil
}

// Start begins consuming messages from Kafka.
func (ec *EventConsumer) Start(ctx context.Context) error {
	if ec.consumer == nil {
		ec.logger.Info("Kafka consumer not configured, skipping start")
		return nil
	}

	ec.logger.Info("Starting event Kafka consumer")
	return ec.consumer.Start(ctx)
}

// Stop gracefully stops the consumer.
func (ec *EventConsumer) Stop() error {
	if ec.consumer == nil {
		return nil
	}

	ec.logger.Info("Stopping event Kafka consumer")
	return ec.consumer.Stop()
}
