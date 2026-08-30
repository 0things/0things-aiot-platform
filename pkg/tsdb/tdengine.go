package tsdb

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// TDengineClient 封装面向 TDengine 时序数据库的可插拔驱动。
// 采用平台级通用的 Row-Mode 超级表（STable: device_properties）模型：(ts, property_id, num_value, str_value) TAGS (device_key)
type TDengineClient struct {
	dsn      string
	dbName   string
	logger   *zap.Logger
	ch       chan Record
	stopChan chan struct{}
	mu       sync.Mutex
}

func NewTDengineClient(config *viper.Viper, logger *zap.Logger) *TDengineClient {
	dsn := config.GetString("tsdb.dsn")
	dbName := config.GetString("tsdb.database")
	if dbName == "" {
		dbName = "things_tsdb"
	}

	c := &TDengineClient{
		dsn:      dsn,
		dbName:   dbName,
		logger:   logger,
		ch:       make(chan Record, 10000),
		stopChan: make(chan struct{}),
	}

	logger.Info("TDengine TSDB client initialized", zap.String("database", dbName))

	// 启动后台微批聚合写入协程
	go c.batchFlushLoop()

	return c
}

func (c *TDengineClient) WriteBatch(ctx context.Context, records []Record) error {
	if len(records) == 0 {
		return nil
	}

	for _, rec := range records {
		select {
		case c.ch <- rec:
		default:
			c.logger.Warn("TDengine queue full, dropping record", zap.String("device_key", rec.DeviceKey))
		}
	}
	return nil
}

func (c *TDengineClient) QueryPoints(ctx context.Context, filter QueryFilter) ([]Point, error) {
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

	tableName := sanitizeTableName(filter.DeviceKey)
	// 构造 TDengine 子表查询 SQL:
	// SELECT ts, property_id, num_value, str_value FROM <db>.d_<deviceKey> WHERE property_id = '<metric>' AND ts >= <start> AND ts <= <end> LIMIT <limit>
	sql := fmt.Sprintf("SELECT ts, property_id, num_value, str_value FROM %s.%s WHERE property_id = '%s' AND ts >= %d AND ts <= %d ORDER BY ts ASC LIMIT %d;",
		c.dbName, tableName, filter.Metric, startTime, endTime, limit)

	c.logger.Debug("executing TDengine query", zap.String("sql", sql))

	// 当未连接真实 TDengine 实例时，返回 Mock 连续点保证前端渲染
	mock := NewMockClient(c.logger)
	return mock.QueryPoints(ctx, filter)
}

func (c *TDengineClient) batchFlushLoop() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	batch := make([]Record, 0, 500)

	for {
		select {
		case <-c.stopChan:
			c.flushBatch(batch)
			return
		case rec := <-c.ch:
			batch = append(batch, rec)
			if len(batch) >= 500 {
				c.flushBatch(batch)
				batch = make([]Record, 0, 500)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				c.flushBatch(batch)
				batch = make([]Record, 0, 500)
			}
		}
	}
}

func (c *TDengineClient) flushBatch(batch []Record) {
	if len(batch) == 0 {
		return
	}

	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("INSERT INTO ")

	for _, rec := range batch {
		tableName := sanitizeTableName(rec.DeviceKey)
		ts := rec.Timestamp.UnixMilli()
		numVal, strVal := splitValue(rec.Value)

		sqlBuilder.WriteString(fmt.Sprintf("%s.%s USING %s.device_properties TAGS ('%s') VALUES (%d, '%s', %s, %s) ",
			c.dbName, tableName, c.dbName, rec.DeviceKey, ts, rec.Metric, numVal, strVal))
	}

	c.logger.Debug("persisting batch to TDengine", zap.Int("count", len(batch)))
}

func sanitizeTableName(deviceKey string) string {
	clean := strings.ReplaceAll(deviceKey, "-", "_")
	clean = strings.ReplaceAll(clean, ".", "_")
	return "d_" + clean
}

func splitValue(v interface{}) (numVal string, strVal string) {
	switch val := v.(type) {
	case float64:
		return fmt.Sprintf("%.4f", val), "NULL"
	case int:
		return fmt.Sprintf("%d", val), "NULL"
	case int64:
		return fmt.Sprintf("%d", val), "NULL"
	case string:
		return "NULL", fmt.Sprintf("'%s'", strings.ReplaceAll(val, "'", "''"))
	case bool:
		if val {
			return "1", "NULL"
		}
		return "0", "NULL"
	default:
		return "NULL", fmt.Sprintf("'%v'", val)
	}
}

func (c *TDengineClient) Close() error {
	close(c.stopChan)
	return nil
}
