package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"aiot-backend/pkg/log"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

var (
	ErrKafkaNotInitialized = errors.New("kafka producer not initialized")
)

type KafkaServiceInterface interface {
	Produce(ctx context.Context, topic string, key, value []byte) error
	ProduceJSON(ctx context.Context, topic string, key string, data any) error
	ProduceAsync(ctx context.Context, topic string, key, value []byte, callback func(error))
	ProduceJSONAsync(ctx context.Context, topic string, key string, data any, callback func(error))
	Flush(ctx context.Context) error
	Close()
}

type KafkaService struct {
	client  *kgo.Client
	logger  *log.Logger
	brokers []string
	enabled bool
}

func NewKafkaService(config *viper.Viper, logger *log.Logger) (*KafkaService, func(), error) {
	brokers := config.GetStringSlice("data.kafka.device.brokers")
	if len(brokers) == 0 {
		logger.Info("Kafka brokers not configured, KafkaService running in disabled mode")
		svc := &KafkaService{
			logger:  logger,
			enabled: false,
		}
		return svc, func() {}, nil
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
	}

	clientID := config.GetString("data.kafka.client_id")
	if clientID != "" {
		opts = append(opts, kgo.ClientID(clientID))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		logger.Error("Failed to create Kafka producer client", zap.Error(err), zap.Strings("brokers", brokers))
		return nil, nil, fmt.Errorf("failed to create kafka client: %w", err)
	}

	// Ping brokers to verify connectivity (warning only, franz-go reconnects automatically)
	if err := client.Ping(context.Background()); err != nil {
		logger.Warn("Failed to ping Kafka brokers (will retry in background)", zap.Error(err), zap.Strings("brokers", brokers))
	} else {
		logger.Info("Kafka producer client initialized successfully", zap.Strings("brokers", brokers))
	}

	svc := &KafkaService{
		client:  client,
		logger:  logger,
		brokers: brokers,
		enabled: true,
	}

	cleanup := func() {
		svc.Close()
	}

	return svc, cleanup, nil
}

func (s *KafkaService) Produce(ctx context.Context, topic string, key, value []byte) error {
	if !s.enabled || s.client == nil {
		return ErrKafkaNotInitialized
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: value,
	}

	results := s.client.ProduceSync(ctx, record)
	if err := results.FirstErr(); err != nil {
		s.logger.Error("Failed to produce Kafka message",
			zap.String("topic", topic),
			zap.Error(err),
		)
		return fmt.Errorf("failed to produce kafka message: %w", err)
	}

	return nil
}

func (s *KafkaService) ProduceJSON(ctx context.Context, topic string, key string, data any) error {
	if !s.enabled || s.client == nil {
		return ErrKafkaNotInitialized
	}

	payload, err := json.Marshal(data)
	if err != nil {
		s.logger.Error("Failed to marshal JSON for Kafka produce",
			zap.String("topic", topic),
			zap.Error(err),
		)
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	var keyBytes []byte
	if key != "" {
		keyBytes = []byte(key)
	}

	return s.Produce(ctx, topic, keyBytes, payload)
}

func (s *KafkaService) ProduceAsync(ctx context.Context, topic string, key, value []byte, callback func(error)) {
	if !s.enabled || s.client == nil {
		if callback != nil {
			callback(ErrKafkaNotInitialized)
		}
		return
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: value,
	}

	s.client.Produce(ctx, record, func(r *kgo.Record, err error) {
		if err != nil {
			s.logger.Error("Failed to produce async Kafka message",
				zap.String("topic", topic),
				zap.Error(err),
			)
		}
		if callback != nil {
			callback(err)
		}
	})
}

func (s *KafkaService) ProduceJSONAsync(ctx context.Context, topic string, key string, data any, callback func(error)) {
	if !s.enabled || s.client == nil {
		if callback != nil {
			callback(ErrKafkaNotInitialized)
		}
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		s.logger.Error("Failed to marshal JSON for async Kafka produce",
			zap.String("topic", topic),
			zap.Error(err),
		)
		if callback != nil {
			callback(fmt.Errorf("failed to marshal json: %w", err))
		}
		return
	}

	var keyBytes []byte
	if key != "" {
		keyBytes = []byte(key)
	}

	s.ProduceAsync(ctx, topic, keyBytes, payload, callback)
}

func (s *KafkaService) Flush(ctx context.Context) error {
	if !s.enabled || s.client == nil {
		return ErrKafkaNotInitialized
	}
	return s.client.Flush(ctx)
}

func (s *KafkaService) Close() {
	if s.client != nil {
		s.client.Close()
		s.logger.Info("Kafka producer client closed")
	}
}

