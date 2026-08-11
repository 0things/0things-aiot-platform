package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"telemetry-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Producer wraps franz-go client for producing messages to Kafka.
type Producer struct {
	client *kgo.Client
	logger *log.Helper
}

// NewProducer creates a new Kafka producer with the given configuration.
func NewProducer(c *conf.Data, logger log.Logger) (*Producer, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "kafka/producer"))

	// Skip initialization if Kafka config is not provided
	if c.Kafka == nil || len(c.Kafka.Brokers) == 0 {
		helper.Info("Kafka configuration not found, skipping producer initialization")
		return nil, func() {}, nil
	}

	helper.Infof("Initializing Kafka producer with brokers: %v", c.Kafka.Brokers)

	// Build franz-go client options
	opts := []kgo.Opt{
		kgo.SeedBrokers(c.Kafka.Brokers...),
	}

	// Add client ID if provided
	if c.Kafka.ClientId != "" {
		opts = append(opts, kgo.ClientID(c.Kafka.ClientId))
	}

	// Add producer-specific options if configured
	if c.Kafka.Producer != nil {
		// Configure batch size for better throughput
		if c.Kafka.Producer.BatchMaxBytes > 0 {
			opts = append(opts, kgo.ProducerBatchMaxBytes(c.Kafka.Producer.BatchMaxBytes))
		} else {
			opts = append(opts, kgo.ProducerBatchMaxBytes(16384)) // 16KB default
		}

		// Configure request retries
		if c.Kafka.Producer.MaxRetries > 0 {
			opts = append(opts, kgo.RequestRetries(int(c.Kafka.Producer.MaxRetries)))
		}

		// Configure request timeout
		if c.Kafka.Producer.Timeout != nil {
			opts = append(opts, kgo.RequestTimeoutOverhead(c.Kafka.Producer.Timeout.AsDuration()))
		}

		// Configure linger for batching
		if c.Kafka.Producer.LingerMs != nil {
			opts = append(opts, kgo.ProducerLinger(c.Kafka.Producer.LingerMs.AsDuration()))
		}

		// Configure compression
		switch c.Kafka.Producer.Compression {
		case "gzip":
			opts = append(opts, kgo.ProducerBatchCompression(kgo.GzipCompression()))
		case "snappy":
			opts = append(opts, kgo.ProducerBatchCompression(kgo.SnappyCompression()))
		case "lz4":
			opts = append(opts, kgo.ProducerBatchCompression(kgo.Lz4Compression()))
		case "zstd":
			opts = append(opts, kgo.ProducerBatchCompression(kgo.ZstdCompression()))
		}

		// Configure required acks
		switch c.Kafka.Producer.RequiredAcks {
		case "none":
			opts = append(opts, kgo.RequiredAcks(kgo.NoAck()))
		case "leader":
			opts = append(opts, kgo.RequiredAcks(kgo.LeaderAck()))
		case "all":
			opts = append(opts, kgo.RequiredAcks(kgo.AllISRAcks()))
		default:
			// Default to leader ack for balance of performance and durability
			opts = append(opts, kgo.RequiredAcks(kgo.LeaderAck()))
		}
	}

	// Create franz-go client
	client, err := kgo.NewClient(opts...)
	if err != nil {
		helper.Errorf("Failed to create Kafka producer client: %v", err)
		return nil, nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// Test connection by pinging brokers
	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		helper.Warnf("Failed to ping Kafka brokers: %v (will retry in background)", err)
		// Don't fail initialization - franz-go will retry connections
	} else {
		helper.Info("Successfully connected to Kafka brokers")
	}

	producer := &Producer{
		client: client,
		logger: helper,
	}

	// Cleanup function
	cleanup := func() {
		helper.Info("Closing Kafka producer")
		client.Close()
	}

	return producer, cleanup, nil
}

// SendMessage sends a message with key-value to the specified topic synchronously.
func (p *Producer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
	if p == nil || p.client == nil {
		return fmt.Errorf("producer not initialized")
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: value,
	}

	// Produce synchronously and wait for result
	results := p.client.ProduceSync(ctx, record)
	if err := results.FirstErr(); err != nil {
		p.logger.Errorf("Failed to produce message to topic %s: %v", topic, err)
		return fmt.Errorf("failed to produce message: %w", err)
	}

	p.logger.Debugf("Produced message to topic %s (partition: %d, offset: %d)",
		topic, results[0].Record.Partition, results[0].Record.Offset)

	return nil
}

// SendJSON encodes data as JSON and sends it to the specified topic synchronously.
func (p *Producer) SendJSON(ctx context.Context, topic string, key string, data interface{}) error {
	if p == nil || p.client == nil {
		return fmt.Errorf("producer not initialized")
	}

	valueBytes, err := json.Marshal(data)
	if err != nil {
		p.logger.Errorf("Failed to marshal JSON: %v", err)
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	return p.SendMessage(ctx, topic, []byte(key), valueBytes)
}

// SendMessageAsync sends a message asynchronously with optional callback.
func (p *Producer) SendMessageAsync(ctx context.Context, topic string, key, value []byte, callback func(error)) {
	if p == nil || p.client == nil {
		if callback != nil {
			callback(fmt.Errorf("producer not initialized"))
		}
		return
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: value,
	}

	// Produce asynchronously with callback
	p.client.Produce(ctx, record, func(r *kgo.Record, err error) {
		if err != nil {
			p.logger.Errorf("Failed to produce async message to topic %s: %v", topic, err)
		} else {
			p.logger.Debugf("Produced async message to topic %s (partition: %d, offset: %d)",
				topic, r.Partition, r.Offset)
		}
		if callback != nil {
			callback(err)
		}
	})
}

// SendJSONAsync encodes data as JSON and sends it to the specified topic asynchronously.
func (p *Producer) SendJSONAsync(ctx context.Context, topic string, key string, data interface{}, callback func(error)) {
	if p == nil || p.client == nil {
		if callback != nil {
			callback(fmt.Errorf("producer not initialized"))
		}
		return
	}

	valueBytes, err := json.Marshal(data)
	if err != nil {
		p.logger.Errorf("Failed to marshal JSON: %v", err)
		if callback != nil {
			callback(fmt.Errorf("failed to marshal json: %w", err))
		}
		return
	}

	p.SendMessageAsync(ctx, topic, []byte(key), valueBytes, callback)
}

// Flush blocks until all buffered records are flushed.
func (p *Producer) Flush(ctx context.Context) error {
	if p == nil || p.client == nil {
		return fmt.Errorf("producer not initialized")
	}

	return p.client.Flush(ctx)
}

// Close gracefully shuts down the producer.
func (p *Producer) Close() {
	if p != nil && p.client != nil {
		p.client.Close()
		p.logger.Info("Kafka producer closed")
	}
}
