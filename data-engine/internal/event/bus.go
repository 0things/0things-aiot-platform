package event

import (
	"fmt"
	"strings"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewBus creates a configured Publisher and Subscriber based on Viper configuration or defaults.
func NewBus(conf *viper.Viper, logger *zap.Logger) (message.Publisher, message.Subscriber, error) {
	driver := "nats"
	if conf != nil {
		if configured := conf.GetString("event.driver"); configured != "" {
			driver = strings.ToLower(configured)
		}
	}

	switch driver {
	case "nats":
		url := "nats://127.0.0.1:4222"
		if conf != nil {
			if configured := conf.GetString("event.nats.url"); configured != "" {
				url = configured
			}
		}
		return NewNATSBus(url, logger)

	case "kafka":
		var brokers []string
		consumerGroup := "0things-consumer"
		if conf != nil {
			brokers = conf.GetStringSlice("event.kafka.brokers")
			if len(brokers) == 0 {
				if configured := conf.GetString("event.kafka.brokers"); configured != "" {
					brokers = []string{configured}
				}
			}
			if configured := conf.GetString("event.kafka.consumer_group"); configured != "" {
				consumerGroup = configured
			}
		}
		if len(brokers) == 0 {
			brokers = []string{"127.0.0.1:9092"}
		}
		return NewKafkaBus(brokers, consumerGroup, logger)

	default:
		return nil, nil, fmt.Errorf("unsupported event driver: %s", driver)
	}
}

// NewNATSBus creates a NATS JetStream Publisher and Subscriber.
func NewNATSBus(url string, logger *zap.Logger) (message.Publisher, message.Subscriber, error) {
	watermillLogger := watermill.NopLogger{}
	marshaler := &nats.JSONMarshaler{}

	pub, err := nats.NewPublisher(nats.PublisherConfig{
		URL:       url,
		Marshaler: marshaler,
		JetStream: nats.JetStreamConfig{
			Disabled:      false,
			AutoProvision: true,
		},
	}, watermillLogger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create NATS publisher: %w", err)
	}

	sub, err := nats.NewSubscriber(nats.SubscriberConfig{
		URL:              url,
		Unmarshaler:      marshaler,
		QueueGroupPrefix: "0things-",
		JetStream: nats.JetStreamConfig{
			Disabled:      false,
			AutoProvision: true,
		},
	}, watermillLogger)
	if err != nil {
		_ = pub.Close()
		return nil, nil, fmt.Errorf("failed to create NATS subscriber: %w", err)
	}

	return pub, sub, nil
}

// NewKafkaBus creates an Apache Kafka Publisher and Subscriber.
func NewKafkaBus(brokers []string, consumerGroup string, logger *zap.Logger) (message.Publisher, message.Subscriber, error) {
	watermillLogger := watermill.NopLogger{}

	pub, err := kafka.NewPublisher(kafka.PublisherConfig{
		Brokers:   brokers,
		Marshaler: kafka.DefaultMarshaler{},
	}, watermillLogger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Kafka publisher: %w", err)
	}

	sub, err := kafka.NewSubscriber(kafka.SubscriberConfig{
		Brokers:       brokers,
		Unmarshaler:   kafka.DefaultMarshaler{},
		ConsumerGroup: consumerGroup,
	}, watermillLogger)
	if err != nil {
		_ = pub.Close()
		return nil, nil, fmt.Errorf("failed to create Kafka subscriber: %w", err)
	}

	return pub, sub, nil
}
