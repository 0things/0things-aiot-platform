package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"telemetry-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

// MessageHandler is the function signature for handling consumed messages.
type MessageHandler func(ctx context.Context, topic string, key, value []byte) error

// TopicConsumer wraps franz-go client for consuming messages from a single Kafka topic.
// Each TopicConsumer has its own client, group, and worker pool for isolation.
type TopicConsumer struct {
	client      *kgo.Client
	logger      *log.Helper
	topic       string
	groupID     string
	handler     MessageHandler
	concurrency int32
	stopCh      chan struct{}
	doneCh      chan struct{}
	running     bool
	mu          sync.Mutex

	// Concurrency control
	workerPool chan struct{}
	wg         sync.WaitGroup
}

// TopicConsumerConfig holds the configuration for creating a TopicConsumer.
type TopicConsumerConfig struct {
	Topic       string
	GroupID     string
	Handler     MessageHandler
	Concurrency int32
}

// NewTopicConsumer creates a new Kafka consumer for a single topic.
func NewTopicConsumer(c *conf.Data, cfg TopicConsumerConfig, logger log.Logger) (*TopicConsumer, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", fmt.Sprintf("kafka/consumer/%s", cfg.Topic)))

	// Skip initialization if Kafka config is not provided
	if c.Kafka == nil || len(c.Kafka.Brokers) == 0 {
		helper.Info("Kafka configuration not found, skipping consumer initialization")
		return nil, func() {}, nil
	}

	// Default concurrency
	concurrency := cfg.Concurrency
	if concurrency <= 0 {
		concurrency = 1
	}

	helper.Infof("Initializing Kafka consumer for topic: %s, group: %s, concurrency: %d",
		cfg.Topic, cfg.GroupID, concurrency)

	// Build franz-go client options
	opts := []kgo.Opt{
		kgo.SeedBrokers(c.Kafka.Brokers...),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeTopics(cfg.Topic),
		kgo.DisableAutoCommit(),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()),
	}

	// Add client ID if provided
	if c.Kafka.ClientId != "" {
		opts = append(opts, kgo.ClientID(c.Kafka.ClientId))
	}

	// Create franz-go client
	client, err := kgo.NewClient(opts...)
	if err != nil {
		helper.Errorf("Failed to create Kafka consumer client: %v", err)
		return nil, nil, fmt.Errorf("failed to create kafka consumer for topic %s: %w", cfg.Topic, err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		helper.Warnf("Failed to ping Kafka brokers: %v (will retry in background)", err)
	} else {
		helper.Info("Successfully connected to Kafka brokers")
	}

	tc := &TopicConsumer{
		client:      client,
		logger:      helper,
		topic:       cfg.Topic,
		groupID:     cfg.GroupID,
		handler:     cfg.Handler,
		concurrency: concurrency,
		stopCh:      make(chan struct{}),
		doneCh:      make(chan struct{}),
		workerPool:  make(chan struct{}, concurrency),
	}

	cleanup := func() {
		helper.Info("Closing Kafka consumer")
		tc.Stop()
	}

	return tc, cleanup, nil
}

// Start begins consuming messages. This is a blocking call.
func (tc *TopicConsumer) Start(ctx context.Context) error {
	if tc == nil || tc.client == nil {
		return nil
	}

	tc.mu.Lock()
	if tc.running {
		tc.mu.Unlock()
		return fmt.Errorf("consumer for topic %s already running", tc.topic)
	}
	tc.running = true
	tc.mu.Unlock()

	defer close(tc.doneCh)

	tc.logger.Infof("Starting Kafka consumer for topic: %s", tc.topic)

	for {
		select {
		case <-ctx.Done():
			tc.logger.Info("Consumer context cancelled, stopping")
			tc.waitForWorkers()
			return ctx.Err()
		case <-tc.stopCh:
			tc.logger.Info("Consumer stop signal received")
			tc.waitForWorkers()
			return nil
		default:
			// Poll for messages with a timeout to check stop signals
			pollCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			fetches := tc.client.PollFetches(pollCtx)
			cancel()

			// Check for errors
			if errs := fetches.Errors(); len(errs) > 0 {
				for _, err := range errs {
					if err.Err == context.DeadlineExceeded || err.Err == context.Canceled {
						continue
					}
					tc.logger.Errorf("Kafka fetch error on topic %s partition %d: %v",
						err.Topic, err.Partition, err.Err)
				}
				continue
			}

			// Process records
			fetches.EachPartition(func(p kgo.FetchTopicPartition) {
				for _, record := range p.Records {
					tc.processRecordAsync(ctx, record)
				}
			})
		}
	}
}

// processRecordAsync processes a record using a worker from the pool.
func (tc *TopicConsumer) processRecordAsync(ctx context.Context, record *kgo.Record) {
	// Acquire worker slot
	select {
	case tc.workerPool <- struct{}{}:
		// Got a slot, process in goroutine
		tc.wg.Add(1)
		go func(rec *kgo.Record) {
			defer tc.wg.Done()
			defer func() { <-tc.workerPool }()

			tc.processRecord(ctx, rec)
		}(record)
	case <-ctx.Done():
		return
	}
}

// processRecord handles a single record.
func (tc *TopicConsumer) processRecord(ctx context.Context, record *kgo.Record) {
	startTime := time.Now()

	// Execute handler
	err := tc.handler(ctx, record.Topic, record.Key, record.Value)

	duration := time.Since(startTime)

	if err != nil {
		tc.logger.Errorf("Handler error for topic %s (partition: %d, offset: %d, duration: %v): %v",
			record.Topic, record.Partition, record.Offset, duration, err)
	} else {
		tc.logger.Debugf("Processed message from topic %s (partition: %d, offset: %d, duration: %v)",
			record.Topic, record.Partition, record.Offset, duration)
	}

	// Mark record for commit
	tc.client.MarkCommitRecords(record)
}

// waitForWorkers waits for all worker goroutines to finish.
func (tc *TopicConsumer) waitForWorkers() {
	done := make(chan struct{})
	go func() {
		tc.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		tc.logger.Info("All workers finished")
	case <-time.After(30 * time.Second):
		tc.logger.Warn("Timeout waiting for workers, some messages may not be processed")
	}
}

// Stop gracefully shuts down the consumer.
func (tc *TopicConsumer) Stop() error {
	if tc == nil {
		return nil
	}

	tc.mu.Lock()
	if !tc.running {
		tc.mu.Unlock()
		return nil
	}
	tc.mu.Unlock()

	tc.logger.Infof("Stopping Kafka consumer for topic: %s", tc.topic)

	// Signal stop
	close(tc.stopCh)

	// Wait for consumer loop to finish with timeout
	select {
	case <-tc.doneCh:
		tc.logger.Info("Consumer polling stopped")
	case <-time.After(30 * time.Second):
		tc.logger.Warn("Consumer stop timeout")
	}

	// Final commit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := tc.client.CommitUncommittedOffsets(ctx); err != nil {
		tc.logger.Errorf("Failed to commit final offsets: %v", err)
	}

	// Close client
	if tc.client != nil {
		tc.client.Close()
		tc.logger.Info("Kafka consumer closed")
	}

	return nil
}

// Topic returns the topic name this consumer is subscribed to.
func (tc *TopicConsumer) Topic() string {
	if tc == nil {
		return ""
	}
	return tc.topic
}
