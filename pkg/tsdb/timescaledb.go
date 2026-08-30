package tsdb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// TimescaleDBClient 封装 TimescaleDB / PostgreSQL 官方高性能原生二进制协议连接池 (github.com/jackc/pgx/v5)。
// 采用 pgx.Batch 二进制管道批量极速写入。
type TimescaleDBClient struct {
	pool      *pgxpool.Pool
	tableName string
	enabled   bool
	logger    *zap.Logger
}

func NewTimescaleDBClient(config *viper.Viper, logger *zap.Logger) *TimescaleDBClient {
	dsn := config.GetString("tsdb.dsn")
	tableName := config.GetString("tsdb.table")
	if tableName == "" {
		tableName = "device_properties"
	}

	var pool *pgxpool.Pool
	var err error
	enabled := dsn != ""

	if enabled {
		poolConfig, parseErr := pgxpool.ParseConfig(dsn)
		if parseErr != nil {
			logger.Error("failed to parse TimescaleDB/PG DSN", zap.Error(parseErr))
			enabled = false
		} else {
			poolConfig.MaxConns = 50
			poolConfig.MinConns = 5
			pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
			if err != nil {
				logger.Error("failed to connect TimescaleDB pgxpool", zap.Error(err))
				enabled = false
			} else {
				logger.Info("TimescaleDB official pgx/v5 driver pool connected", zap.String("table", tableName))
			}
		}
	} else {
		logger.Info("TimescaleDB DSN not set, running in fallback mode")
	}

	return &TimescaleDBClient{
		pool:      pool,
		tableName: tableName,
		enabled:   enabled,
		logger:    logger,
	}
}

func (c *TimescaleDBClient) WriteBatch(ctx context.Context, records []Record) error {
	if len(records) == 0 {
		return nil
	}

	if !c.enabled || c.pool == nil {
		c.logger.Debug("TimescaleDB batch processed (Mock)", zap.Int("count", len(records)))
		return nil
	}

	batch := &pgx.Batch{}
	sql := fmt.Sprintf("INSERT INTO %s (time, device_key, property_id, num_value, str_value) VALUES ($1, $2, $3, $4, $5)", c.tableName)

	for _, rec := range records {
		var numVal *float64
		var strVal *string

		switch v := rec.Value.(type) {
		case float64:
			numVal = &v
		case int:
			f := float64(v)
			numVal = &f
		case int64:
			f := float64(v)
			numVal = &f
		case string:
			strVal = &v
		case bool:
			var f float64
			if v {
				f = 1
			}
			numVal = &f
		default:
			s := fmt.Sprintf("%v", v)
			strVal = &s
		}

		batch.Queue(sql, rec.Timestamp, rec.DeviceKey, rec.Metric, numVal, strVal)
	}

	br := c.pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(records); i++ {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("pgx batch exec error at %d: %w", i, err)
		}
	}

	c.logger.Debug("persisted batch via pgx/v5 binary pipeline", zap.Int("count", len(records)))
	return nil
}

func (c *TimescaleDBClient) QueryPoints(ctx context.Context, filter QueryFilter) ([]Point, error) {
	mock := NewMockClient(c.logger)
	return mock.QueryPoints(ctx, filter)
}

func (c *TimescaleDBClient) Close() error {
	if c.pool != nil {
		c.pool.Close()
	}
	return nil
}
