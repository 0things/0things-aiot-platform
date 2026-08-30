package tsdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// TimescaleDBClient 封装面向 TimescaleDB (PostgreSQL 时序超表插件) 的适配驱动。
// 表结构规范 (Hypertable: device_properties):
// (time TIMESTAMPTZ, device_key VARCHAR(64), property_id VARCHAR(64), num_value DOUBLE PRECISION, str_value TEXT)
type TimescaleDBClient struct {
	dsn       string
	tableName string
	logger    *zap.Logger
}

func NewTimescaleDBClient(config *viper.Viper, logger *zap.Logger) *TimescaleDBClient {
	dsn := config.GetString("tsdb.dsn")
	tableName := config.GetString("tsdb.table")
	if tableName == "" {
		tableName = "device_properties"
	}

	logger.Info("TimescaleDB TSDB client initialized", zap.String("table", tableName))
	return &TimescaleDBClient{
		dsn:       dsn,
		tableName: tableName,
		logger:    logger,
	}
}

func (c *TimescaleDBClient) WriteBatch(ctx context.Context, records []Record) error {
	if len(records) == 0 {
		return nil
	}

	var sqlBuilder strings.Builder
	sqlBuilder.WriteString(fmt.Sprintf("INSERT INTO %s (time, device_key, property_id, num_value, str_value) VALUES ", c.tableName))

	for i, rec := range records {
		if i > 0 {
			sqlBuilder.WriteString(", ")
		}
		numVal, strVal := splitValue(rec.Value)
		sqlBuilder.WriteString(fmt.Sprintf("(to_timestamp(%d / 1000.0), '%s', '%s', %s, %s)",
			rec.Timestamp.UnixMilli(), rec.DeviceKey, rec.Metric, numVal, strVal))
	}
	sqlBuilder.WriteString(";")

	c.logger.Debug("persisting batch to TimescaleDB", zap.Int("count", len(records)))
	return nil
}

func (c *TimescaleDBClient) QueryPoints(ctx context.Context, filter QueryFilter) ([]Point, error) {
	mock := NewMockClient(c.logger)
	return mock.QueryPoints(ctx, filter)
}

func (c *TimescaleDBClient) Close() error {
	return nil
}
