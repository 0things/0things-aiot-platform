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

const defaultTelemetryConcurrency = 10

// TelemetryConsumer handles consuming telemetry events from Kafka.
type TelemetryConsumer struct {
	uc       *biz.TelemetryUsecase
	consumer *kafka.TopicConsumer
	logger   *log.Helper
}

// NewTelemetryConsumer creates a new telemetry consumer.
func NewTelemetryConsumer(uc *biz.TelemetryUsecase, c *conf.Data, logger log.Logger) (*TelemetryConsumer, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "consumer/telemetry"))

	tc := &TelemetryConsumer{
		uc:     uc,
		logger: helper,
	}

	// Create topic consumer
	consumer, cleanup, err := kafka.NewTopicConsumer(c, kafka.TopicConsumerConfig{
		Topic:       kafka.TopicTelemetry,
		GroupID:     kafka.GroupTelemetry,
		Handler:     tc.handleTelemetryEvent,
		Concurrency: defaultTelemetryConcurrency,
	}, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create telemetry consumer: %w", err)
	}

	tc.consumer = consumer

	if consumer != nil {
		helper.Infof("TelemetryConsumer initialized for topic: %s, group: %s", kafka.TopicTelemetry, kafka.GroupTelemetry)
	} else {
		helper.Info("Kafka not configured, TelemetryConsumer running in no-op mode")
	}

	return tc, cleanup, nil
}

// handleTelemetryEvent processes telemetry events from Kafka.
func (tc *TelemetryConsumer) handleTelemetryEvent(ctx context.Context, topic string, key, value []byte) error {
	tc.logger.Debugf("Received message from topic %s with key: %s", topic, string(key))

	// Step 1: Deserialize JSON
	var event biz.TelemetryEvent
	if err := json.Unmarshal(value, &event); err != nil {
		tc.logger.Errorf("Failed to unmarshal telemetry event: %v", err)
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Step 2: Basic validation
	if event.Type == "" {
		tc.logger.Error("Event type is required")
		return fmt.Errorf("type is required")
	}
	if event.Timestamp == 0 {
		tc.logger.Error("Timestamp is required")
		return fmt.Errorf("timestamp is required")
	}

	// Step 3: Pass to business logic layer
	if err := tc.uc.ProcessEvent(ctx, &event); err != nil {
		tc.logger.Errorf("Failed to process telemetry event: %v", err)
		return err
	}

	return nil
}

// Start begins consuming messages from Kafka.
func (tc *TelemetryConsumer) Start(ctx context.Context) error {
	if tc.consumer == nil {
		tc.logger.Info("Kafka consumer not configured, skipping start")
		return nil
	}

	tc.logger.Info("Starting telemetry Kafka consumer")
	return tc.consumer.Start(ctx)
}

// Stop gracefully stops the consumer.
func (tc *TelemetryConsumer) Stop() error {
	if tc.consumer == nil {
		return nil
	}

	tc.logger.Info("Stopping telemetry Kafka consumer")
	return tc.consumer.Stop()
}
