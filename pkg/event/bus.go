package event

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	_ "github.com/glebarez/go-sqlite"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewBus creates a configured Publisher and Subscriber based on Viper configuration or defaults.
func NewBus(conf *viper.Viper, logger *zap.Logger) (message.Publisher, message.Subscriber, error) {
	driver := "nats"
	if conf != nil {
		if d := conf.GetString("event.driver"); d != "" {
			driver = strings.ToLower(d)
		}
	}

	switch driver {
	case "memory":
		return NewMemoryBus()

	case "nats":
		natsURL := "nats://127.0.0.1:4222"
		if conf != nil {
			if u := conf.GetString("event.nats.url"); u != "" {
				natsURL = u
			}
		}
		return NewNATSBus(natsURL, logger)

	case "kafka":
		var brokers []string
		consumerGroup := "0things-consumer"
		if conf != nil {
			brokers = conf.GetStringSlice("event.kafka.brokers")
			if len(brokers) == 0 {
				if b := conf.GetString("event.kafka.brokers"); b != "" {
					brokers = []string{b}
				}
			}
			if cg := conf.GetString("event.kafka.consumer_group"); cg != "" {
				consumerGroup = cg
			}
		}
		if len(brokers) == 0 {
			brokers = []string{"127.0.0.1:9092"}
		}
		return NewKafkaBus(brokers, consumerGroup, logger)

	case "mysql":
		dsn := "root:123456@tcp(127.0.0.1:3306)/0things?parseTime=true"
		consumerGroup := "0things-consumer"
		if conf != nil {
			if d := conf.GetString("event.mysql.dsn"); d != "" {
				dsn = d
			}
			if cg := conf.GetString("event.mysql.consumer_group"); cg != "" {
				consumerGroup = cg
			}
		}
		return NewMySQLBus(dsn, consumerGroup, logger)

	case "postgres":
		dsn := "postgres://postgres:123456@127.0.0.1:5432/0things?sslmode=disable"
		consumerGroup := "0things-consumer"
		if conf != nil {
			if d := conf.GetString("event.postgres.dsn"); d != "" {
				dsn = d
			}
			if cg := conf.GetString("event.postgres.consumer_group"); cg != "" {
				consumerGroup = cg
			}
		}
		return NewPostgresBus(dsn, consumerGroup, logger)

	case "sqlite":
		dsn := "storage/event_bus.db"
		consumerGroup := "0things-consumer"
		if conf != nil {
			if path := conf.GetString("event.sqlite.path"); path != "" {
				dsn = path
			} else if d := conf.GetString("event.sqlite.dsn"); d != "" {
				dsn = d
			}
			if cg := conf.GetString("event.sqlite.consumer_group"); cg != "" {
				consumerGroup = cg
			}
		}
		return NewSQLiteBus(dsn, consumerGroup, logger)

	default:
		return nil, nil, fmt.Errorf("unsupported event driver: %s", driver)
	}
}

// NewMemoryBus creates an in-memory GoChannel Pub/Sub for testing and single-process setups.
func NewMemoryBus() (message.Publisher, message.Subscriber, error) {
	pubSub := gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer:            1024,
		BlockPublishUntilSubscriberAck: true,
	}, watermill.NopLogger{})
	return pubSub, pubSub, nil
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
		pub.Close()
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
		pub.Close()
		return nil, nil, fmt.Errorf("failed to create Kafka subscriber: %w", err)
	}

	return pub, sub, nil
}

// NewMySQLBus creates a MySQL backed Pub/Sub.
func NewMySQLBus(dsn string, consumerGroup string, logger *zap.Logger) (message.Publisher, message.Subscriber, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}
	return NewSQLBusWithDB(db, "mysql", consumerGroup, logger)
}

// NewPostgresBus creates a PostgreSQL backed Pub/Sub.
func NewPostgresBus(dsn string, consumerGroup string, logger *zap.Logger) (message.Publisher, message.Subscriber, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}
	return NewSQLBusWithDB(db, "postgres", consumerGroup, logger)
}

// NewSQLiteBus creates a SQLite backed Pub/Sub.
func NewSQLiteBus(dsn string, consumerGroup string, logger *zap.Logger) (message.Publisher, message.Subscriber, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open SQLite connection: %w", err)
	}
	return NewSQLBusWithDB(db, "sqlite", consumerGroup, logger)
}

// NewSQLBusWithDB initializes a Watermill SQL Publisher and Subscriber on an existing database connection.
func NewSQLBusWithDB(db *sql.DB, dbType string, consumerGroup string, logger *zap.Logger) (message.Publisher, message.Subscriber, error) {
	watermillLogger := watermill.NopLogger{}

	var schemaAdapter watermillSQL.SchemaAdapter
	var offsetsAdapter watermillSQL.OffsetsAdapter

	switch strings.ToLower(dbType) {
	case "postgres":
		schemaAdapter = watermillSQL.DefaultPostgreSQLSchema{}
		offsetsAdapter = watermillSQL.DefaultPostgreSQLOffsetsAdapter{}
	case "sqlite":
		schemaAdapter = SQLiteSchema{}
		offsetsAdapter = SQLiteOffsetsAdapter{}
	case "mysql":
		schemaAdapter = watermillSQL.DefaultMySQLSchema{}
		offsetsAdapter = watermillSQL.DefaultMySQLOffsetsAdapter{}
	default:
		return nil, nil, fmt.Errorf("unsupported SQL event driver: %s", dbType)
	}

	pub, err := watermillSQL.NewPublisher(
		db,
		watermillSQL.PublisherConfig{
			SchemaAdapter:        schemaAdapter,
			AutoInitializeSchema: true,
		},
		watermillLogger,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SQL publisher: %w", err)
	}

	sub, err := watermillSQL.NewSubscriber(
		db,
		watermillSQL.SubscriberConfig{
			ConsumerGroup:    consumerGroup,
			SchemaAdapter:    schemaAdapter,
			OffsetsAdapter:   offsetsAdapter,
			InitializeSchema: true,
			PollInterval:     100 * time.Millisecond,
			ResendInterval:   500 * time.Millisecond,
			RetryInterval:    500 * time.Millisecond,
		},
		watermillLogger,
	)
	if err != nil {
		pub.Close()
		return nil, nil, fmt.Errorf("failed to create SQL subscriber: %w", err)
	}

	return pub, sub, nil
}
