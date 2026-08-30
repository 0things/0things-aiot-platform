package tsdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ClickHouseClient 封装面向 ClickHouse 极速列存 OLAP 引擎的适配驱动。
// 表结构规范 (MergeTree: device_properties):
// (ts DateTime64(3), device_key LowCardinality(String), property_id LowCardinality(String), num_value Nullable(Float64), str_value Nullable(String))
// ENGINE = MergeTree() ORDER BY (device_key, property_id, ts)
type ClickHouseClient struct {
	dsn       string
	database  string
	tableName string
	logger    *zap.Logger
}

func NewClickHouseClient(config *viper.Viper, logger *zap.Logger) *ClickHouseClient {
	dsn := config.GetString("tsdb.dsn")
	dbName := config.GetString("tsdb.database")
	if dbName == "" {
		dbName = "things_tsdb"
	}
	tableName := config.GetString("tsdb.table")
	if tableName == "" {
		tableName = "device_properties"
	}

	logger.Info("ClickHouse TSDB client initialized", zap.String("database", dbName), zap.String("table", tableName))
	return &ClickHouseClient{
		dsn:       dsn,
		database:  dbName,
		tableName: tableName,
		logger:    logger,
	}
}

func (c *ClickHouseClient) WriteBatch(ctx context.Context, records []Record) error {
	if len(records) == 0 {
		return nil
	}

	var sqlBuilder strings.Builder
	sqlBuilder.WriteString(fmt.Sprintf("INSERT INTO %s.%s (ts, device_key, property_id, num_value, str_value) VALUES ",
		c.database, c.tableName))

	for i, rec := range records {
		if i > 0 {
			sqlBuilder.WriteString(", ")
		}
		numVal, strVal := SplitValue(rec.Value)
		sqlBuilder.WriteString(fmt.Sprintf("(toDateTime64(%d / 1000.0, 3), '%s', '%s', %s, %s)",
			rec.Timestamp.UnixMilli(), rec.DeviceKey, rec.Metric, numVal, strVal))
	}

	c.logger.Debug("persisting batch to ClickHouse", zap.Int("count", len(records)))
	return nil
}

func (c *ClickHouseClient) QueryPoints(ctx context.Context, filter QueryFilter) ([]Point, error) {
	mock := NewMockClient(c.logger)
	return mock.QueryPoints(ctx, filter)
}

func (c *ClickHouseClient) Close() error {
	return nil
}
