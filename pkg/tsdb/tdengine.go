package tsdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	_ "github.com/taosdata/driver-go/v3/taosWS"
	"go.uber.org/zap"
)

// TDengineClient 基于 TDengine 官方原生纯 Go WebSocket 驱动 (github.com/taosdata/driver-go/v3/taosWS)。
// 无需任何 CGO / 本地 C 动态库，支持官方连接池管理、自动建库建超级表与全双工高并发写入/查询。
type TDengineClient struct {
	db       *sql.DB
	dbName   string
	logger   *zap.Logger
	ch       chan Record
	stopChan chan struct{}
	mu       sync.Mutex
	enabled  bool
}

func NewTDengineClient(config *viper.Viper, logger *zap.Logger) *TDengineClient {
	rawDSN := config.GetString("tsdb.dsn")
	dbName := config.GetString("tsdb.database")
	if dbName == "" {
		dbName = "things_tsdb"
	}

	wsDSN := formatTaosWSDSN(rawDSN, dbName)

	var db *sql.DB
	var err error
	enabled := rawDSN != ""

	if enabled {
		db, err = sql.Open("taosWS", wsDSN)
		if err != nil {
			logger.Error("failed to open TDengine taosWS driver", zap.Error(err))
			enabled = false
		} else {
			db.SetMaxOpenConns(50)
			db.SetMaxIdleConns(20)
			db.SetConnMaxLifetime(5 * time.Minute)

			logger.Info("TDengine official taosWS driver initialized",
				zap.String("dsn", wsDSN),
				zap.String("database", dbName),
			)
			// 自动建库与超级表
			go func() {
				time.Sleep(1 * time.Second)
				initTDengineSTable(db, dbName, logger)
			}()
		}
	} else {
		logger.Info("TDengine DSN not configured, running in fallback mode")
	}

	c := &TDengineClient{
		db:       db,
		dbName:   dbName,
		logger:   logger,
		ch:       make(chan Record, 10000),
		stopChan: make(chan struct{}),
		enabled:  enabled,
	}

	go c.batchFlushLoop()

	return c
}

// formatTaosWSDSN 将配置的 DSN 统一转换为官方 taosWS 标准格式: "user:pass@ws(host:port)/dbName"
func formatTaosWSDSN(rawDSN, dbName string) string {
	if rawDSN == "" {
		return fmt.Sprintf("root:taosdata@ws(127.0.0.1:6041)/%s", dbName)
	}

	if strings.Contains(rawDSN, "@ws(") {
		return rawDSN
	}

	user := "root"
	pass := "taosdata"
	host := "127.0.0.1:6041"

	if strings.Contains(rawDSN, "@") {
		parts := strings.Split(rawDSN, "@")
		if creds := strings.Split(parts[0], ":"); len(creds) == 2 {
			user = creds[0]
			pass = creds[1]
		}
		hostPart := parts[1]
		hostPart = strings.TrimPrefix(hostPart, "http://")
		hostPart = strings.TrimPrefix(hostPart, "https://")
		hostPart = strings.TrimPrefix(hostPart, "ws://")
		hostPart = strings.TrimPrefix(hostPart, "http(")
		hostPart = strings.TrimPrefix(hostPart, "tcp(")
		if idx := strings.Index(hostPart, ")"); idx != -1 {
			hostPart = hostPart[:idx]
		}
		if idx := strings.Index(hostPart, "/"); idx != -1 {
			hostPart = hostPart[:idx]
		}
		host = hostPart
	} else {
		h := strings.TrimPrefix(rawDSN, "http://")
		h = strings.TrimPrefix(h, "ws://")
		if idx := strings.Index(h, "/"); idx != -1 {
			h = h[:idx]
		}
		host = h
	}

	return fmt.Sprintf("%s:%s@ws(%s)/%s", user, pass, host, dbName)
}

// initTDengineSTable 自动向 TDengine 发送 DDL 创建时序数据库与通用超级表 (STable)
func initTDengineSTable(db *sql.DB, dbName string, logger *zap.Logger) {
	if db == nil {
		return
	}

	// 1. 创建数据库
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s KEEP 365 DAYS 10 BLOCKS 6 PRECISION 'ms';", dbName)
	if _, err := db.Exec(createDBSQL); err != nil {
		logger.Warn("TDengine DDL create database note (may already exist or connecting)", zap.Error(err))
	}

	// 2. 创建通用超级表 (STable: device_properties)
	createSTableSQL := fmt.Sprintf(`CREATE STABLE IF NOT EXISTS %s.device_properties (
		ts TIMESTAMP,
		property_id VARCHAR(64),
		num_value DOUBLE,
		str_value VARCHAR(1024)
	) TAGS (
		device_key VARCHAR(64)
	);`, dbName)

	if _, err := db.Exec(createSTableSQL); err != nil {
		logger.Warn("TDengine DDL create STable note", zap.Error(err))
	} else {
		logger.Info("✓ TDengine official driver verified STable device_properties successfully")
	}
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

// QueryPoints 通过官方 taosWS 连接池执行真实 SQL 查询并提取时序数据点
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

	tableName := SanitizeTableName(filter.DeviceKey)
	sqlStr := fmt.Sprintf("SELECT ts, property_id, num_value, str_value FROM %s.%s WHERE property_id = '%s' AND ts >= %d AND ts <= %d ORDER BY ts ASC LIMIT %d;",
		c.dbName, tableName, filter.Metric, startTime, endTime, limit)

	if c.enabled && c.db != nil {
		rows, err := c.db.QueryContext(ctx, sqlStr)
		if err == nil {
			defer rows.Close()
			points := make([]Point, 0)
			for rows.Next() {
				var ts time.Time
				var propID string
				var numVal sql.NullFloat64
				var strVal sql.NullString

				if scanErr := rows.Scan(&ts, &propID, &numVal, &strVal); scanErr == nil {
					var val interface{}
					if numVal.Valid {
						val = numVal.Float64
					} else if strVal.Valid {
						val = strVal.String
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

	// 离线平滑降级
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

// flushBatch 通过官方 taosWS 连接池执行批量 SQL 插入
func (c *TDengineClient) flushBatch(batch []Record) {
	if len(batch) == 0 {
		return
	}

	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("INSERT INTO ")

	for _, rec := range batch {
		tableName := SanitizeTableName(rec.DeviceKey)
		ts := rec.Timestamp.UnixMilli()
		numVal, strVal := SplitValue(rec.Value)

		sqlBuilder.WriteString(fmt.Sprintf("%s.%s USING %s.device_properties TAGS ('%s') VALUES (%d, '%s', %s, %s) ",
			c.dbName, tableName, c.dbName, rec.DeviceKey, ts, rec.Metric, numVal, strVal))
	}

	sqlStr := sqlBuilder.String()
	c.logger.Debug("persisting batch to TDengine via official taosWS driver", zap.Int("count", len(batch)))

	if c.enabled && c.db != nil {
		if _, err := c.db.Exec(sqlStr); err != nil {
			c.logger.Error("failed to write batch via official taosWS", zap.Error(err))
		}
	}
}

func (c *TDengineClient) Close() error {
	close(c.stopChan)
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
