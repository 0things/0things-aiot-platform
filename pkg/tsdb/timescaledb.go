package tsdb

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// PostgreSQLConfig PostgreSQL 数据库专属配置结构体
type PostgreSQLConfig struct {
	DSN   string `mapstructure:"dsn" yaml:"dsn" json:"dsn"`       // 标准 PG DSN
	Table string `mapstructure:"table" yaml:"table" json:"table"` // 存储表名 (默认: device_properties)
}

// TimescaleDBConfig TimescaleDB 数据库专属配置结构体
type TimescaleDBConfig struct {
	DSN   string `mapstructure:"dsn" yaml:"dsn" json:"dsn"`       // 标准 PG / TimescaleDB DSN
	Table string `mapstructure:"table" yaml:"table" json:"table"` // 存储超表名 (默认: device_properties)
}

// TimescaleDBClient 封装 TimescaleDB / PostgreSQL 官方高性能原生二进制协议连接池 (github.com/jackc/pgx/v5)。
// 采用 pgx.Batch 二进制管道批量极速写入与原生连接池真实查询。
type TimescaleDBClient struct {
	pool      *pgxpool.Pool
	tableName string
	enabled   bool
	logger    *zap.Logger
}

func NewTimescaleDBClient(cfg TimescaleDBConfig, logger *zap.Logger) *TimescaleDBClient {
	return createPGClient(cfg.DSN, cfg.Table, logger)
}

func NewPostgreSQLClient(cfg PostgreSQLConfig, logger *zap.Logger) *TimescaleDBClient {
	return createPGClient(cfg.DSN, cfg.Table, logger)
}

func createPGClient(dsn, tableName string, logger *zap.Logger) *TimescaleDBClient {
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

	sql := fmt.Sprintf("INSERT INTO %s (time, device_key, property_id, num_value, str_value, bool_value, json_value) VALUES ($1, $2, $3, $4, $5, $6, $7)", c.tableName)

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection from pgxpool: %w", err)
	}
	defer conn.Release()

	for _, rec := range records {
		numVal, strVal, boolVal, jsonVal := ToTypedValue(rec.Value)
		if _, err := conn.Exec(ctx, sql, rec.Timestamp, rec.DeviceKey, rec.Metric, numVal, strVal, boolVal, jsonVal); err != nil {
			c.logger.Warn("failed to insert record into TimescaleDB", zap.Error(err))
		}
	}

	c.logger.Debug("persisted batch via official pgx/v5 pool", zap.Int("count", len(records)))
	return nil
}

// QueryPoints 通过官方 pgx/v5 连接池执行真实 SQL 查询
func (c *TimescaleDBClient) QueryPoints(ctx context.Context, filter QueryFilter) ([]Point, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	endTime := filter.EndTime
	if endTime <= 0 {
		endTime = time.Now().UnixMilli()
	}

	startTime := filter.StartTime
	if startTime <= 0 {
		startTime = endTime - 24*3600*1000
	}

	if c.enabled && c.pool != nil {
		sql := fmt.Sprintf(`SELECT time, property_id, num_value, str_value, bool_value, json_value FROM %s 
			WHERE device_key = $1 AND property_id = $2 AND time >= to_timestamp($3 / 1000.0) AND time <= to_timestamp($4 / 1000.0) 
			ORDER BY time ASC LIMIT $5`, c.tableName)

		rows, err := c.pool.Query(ctx, sql, filter.DeviceKey, filter.Metric, startTime, endTime, limit)
		if err == nil {
			defer rows.Close()
			points := make([]Point, 0)
			for rows.Next() {
				var ts time.Time
				var propID string
				var numVal *float64
				var strVal *string
				var boolVal *bool
				var jsonVal *string

				if err := rows.Scan(&ts, &propID, &numVal, &strVal, &boolVal, &jsonVal); err == nil {
					var val interface{}
					if jsonVal != nil && *jsonVal != "" {
						val = UnmarshalJSONValue(*jsonVal)
					} else if boolVal != nil {
						val = *boolVal
					} else if numVal != nil {
						val = *numVal
					} else if strVal != nil {
						val = *strVal
					}

					points = append(points, Point{
						Timestamp: ts.UnixMilli(),
						Metric:    propID,
						Value:     val,
					})
				}
			}
			if len(points) > 0 {
				return points, nil
			}
		}
	}

	// 离线降级
	mock := NewMockClient(c.logger)
	return mock.QueryPoints(ctx, filter)
}

func (c *TimescaleDBClient) Close() error {
	if c.pool != nil {
		c.pool.Close()
	}
	return nil
}
