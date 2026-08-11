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

// ClientConnectState needs serial processing to maintain order
const defaultClientConnectStateConcurrency = 1

// ClientConnectStateConsumer handles consuming client connect state events from Kafka.
type ClientConnectStateConsumer struct {
	uc       *biz.ClientConnectStateUsecase
	consumer *kafka.TopicConsumer
	logger   *log.Helper
}

// NewClientConnectStateConsumer creates a new client connect state consumer.
func NewClientConnectStateConsumer(uc *biz.ClientConnectStateUsecase, c *conf.Data, logger log.Logger) (*ClientConnectStateConsumer, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "consumer/client_connect_state"))

	ccsc := &ClientConnectStateConsumer{
		uc:     uc,
		logger: helper,
	}

	// Create topic consumer with concurrency=1 for serial processing
	consumer, cleanup, err := kafka.NewTopicConsumer(c, kafka.TopicConsumerConfig{
		Topic:       kafka.TopicClientConnectState,
		GroupID:     kafka.GroupClientConnectState,
		Handler:     ccsc.handleEvent,
		Concurrency: defaultClientConnectStateConcurrency,
	}, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client connect state consumer: %w", err)
	}

	ccsc.consumer = consumer

	if consumer != nil {
		helper.Infof("ClientConnectStateConsumer initialized for topic: %s, group: %s (serial mode)",
			kafka.TopicClientConnectState, kafka.GroupClientConnectState)
	} else {
		helper.Info("Kafka not configured, ClientConnectStateConsumer running in no-op mode")
	}

	return ccsc, cleanup, nil
}

// handleEvent processes client connect state events from Kafka.
func (ccsc *ClientConnectStateConsumer) handleEvent(ctx context.Context, topic string, key, value []byte) error {
	ccsc.logger.Debugf("Received event from topic %s with key: %s", topic, string(key))

	// Step 1: Deserialize JSON
	var event biz.ClientConnectStateEvent
	if err := json.Unmarshal(value, &event); err != nil {
		ccsc.logger.Errorf("Failed to unmarshal client connect state event: %v", err)
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Step 2: Basic validation
	if event.ClientID == "" {
		ccsc.logger.Error("ClientID is required")
		return fmt.Errorf("clientid is required")
	}
	if event.Event == "" {
		ccsc.logger.Error("Event is required")
		return fmt.Errorf("event is required")
	}

	ccsc.logger.Infof("Received client connect state event: event=%s clientid=%s device_key=%s timestamp=%d reason=%s",
		event.Event, event.ClientID, event.DeviceKey, event.Timestamp, event.Reason)

	// Step 3: Pass to business logic layer
	if err := ccsc.uc.ProcessEvent(ctx, &event); err != nil {
		ccsc.logger.Errorf("Failed to process client connect state event: %v", err)
		return err
	}

	return nil
}

// Start begins consuming messages from Kafka.
func (ccsc *ClientConnectStateConsumer) Start(ctx context.Context) error {
	if ccsc.consumer == nil {
		ccsc.logger.Info("Kafka consumer not configured, skipping start")
		return nil
	}

	ccsc.logger.Info("Starting client connect state Kafka consumer")
	return ccsc.consumer.Start(ctx)
}

// Stop gracefully stops the consumer.
func (ccsc *ClientConnectStateConsumer) Stop() error {
	if ccsc.consumer == nil {
		return nil
	}

	ccsc.logger.Info("Stopping client connect state Kafka consumer")
	return ccsc.consumer.Stop()
}
