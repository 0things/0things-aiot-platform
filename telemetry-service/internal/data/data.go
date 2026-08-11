package data

import (
	"context"
	"telemetry-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewGreeterRepo,
	NewTelemetryEventRepo,
	NewEventRepo,
	NewClientConnectStateRepo,
)

// Data .
type Data struct {
	rdb *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(logger)

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Username:     c.Redis.Username,
		Password:     c.Redis.Password,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		DB:           int(c.Redis.Db),
	})

	// Test Redis connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		l.Errorf("Failed to connect to Redis: %v", err)
		return nil, nil, err
	}
	l.Info("Successfully connected to Redis")

	cleanup := func() {
		l.Info("closing the data resources")
		if err := rdb.Close(); err != nil {
			l.Errorf("Failed to close Redis: %v", err)
		}
	}
	return &Data{
		rdb: rdb,
	}, cleanup, nil
}
