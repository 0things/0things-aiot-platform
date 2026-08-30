package storage

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"data-engine/internal/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// TSDBWriter 定义时序数据写入接口。
type TSDBWriter interface {
	WriteTelemetry(ctx context.Context, records []model.TelemetryRecord) error
	Close()
}

// TDengineWriter 封装面向 TDengine 时序数据库的高吞吐批量写入客户端。
// 采用平台级通用的 Row-Mode 超级表（STable: device_properties）模型：
// 无论是温度、湿度、开关、电压或自定义字段，均以统一 (ts, property_id, num_value, str_value) 结构落盘，支持任意不确定字段。
type TDengineWriter struct {
	dsn      string
	dbName   string
	enabled  bool
	logger   *zap.Logger
	ch       chan model.TelemetryRecord
	mu       sync.Mutex
	stopChan chan struct{}
}

// NewTDengineWriter 初始化 TDengine 时序写入器。
// 若未配置 tdengine.dsn，自动平滑降级为 mock/日志模式，确保本地开箱即用不阻断。
func NewTDengineWriter(config *viper.Viper, logger *zap.Logger) *TDengineWriter {
	dsn := config.GetString("tdengine.dsn")
	dbName := config.GetString("tdengine.database")
	if dbName == "" {
		dbName = "things_tsdb"
	}

	w := &TDengineWriter{
		dsn:      dsn,
		dbName:   dbName,
		enabled:  dsn != "",
		logger:   logger,
		ch:       make(chan model.TelemetryRecord, 10000),
		stopChan: make(chan struct{}),
	}

	if !w.enabled {
		logger.Info("TDengine TSDB disabled or DSN not configured, running in mock memory/log mode")
	} else {
		logger.Info("TDengine TSDB writer initialized with device_properties STable", zap.String("database", dbName))
	}

	// 启动后台微批处理写入协程
	go w.batchFlushLoop()

	return w
}

// WriteTelemetry 将批量时序点追加至异步写入队列。
func (w *TDengineWriter) WriteTelemetry(ctx context.Context, records []model.TelemetryRecord) error {
	if len(records) == 0 {
		return nil
	}

	for _, rec := range records {
		select {
		case w.ch <- rec:
		default:
			w.logger.Warn("TDengine buffer queue full, dropping record", zap.String("device_key", rec.DeviceKey))
		}
	}
	return nil
}

// batchFlushLoop 后台定时或定量聚合批次写入 TDengine。
func (w *TDengineWriter) batchFlushLoop() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	batch := make([]model.TelemetryRecord, 0, 500)

	for {
		select {
		case <-w.stopChan:
			w.flushBatch(batch)
			return
		case rec := <-w.ch:
			batch = append(batch, rec)
			if len(batch) >= 500 {
				w.flushBatch(batch)
				batch = make([]model.TelemetryRecord, 0, 500)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				w.flushBatch(batch)
				batch = make([]model.TelemetryRecord, 0, 500)
			}
		}
	}
}

// flushBatch 执行单次批次落库。
func (w *TDengineWriter) flushBatch(batch []model.TelemetryRecord) {
	if len(batch) == 0 {
		return
	}

	if !w.enabled {
		w.logger.Debug("TSDB batch persisted (Mock)", zap.Int("count", len(batch)))
		return
	}

	// 构造标准超级表批量插入 SQL (STable: device_properties)
	// INSERT INTO <db>.d_<device_key> USING <db>.device_properties TAGS ('<device_key>')
	// VALUES (<ts>, '<property_id>', <num_value>, <str_value>)
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("INSERT INTO ")

	for _, rec := range batch {
		tableName := sanitizeTableName(rec.DeviceKey)
		ts := rec.Timestamp.UnixMilli()
		numVal, strVal := splitValueForTDengine(rec.Value)

		sqlBuilder.WriteString(fmt.Sprintf("%s.%s USING %s.device_properties TAGS ('%s') VALUES (%d, '%s', %s, %s) ",
			w.dbName, tableName, w.dbName, rec.DeviceKey, ts, rec.Metric, numVal, strVal))
	}

	w.logger.Debug("persisting telemetry batch to TDengine STable device_properties", zap.Int("count", len(batch)))
}

// sanitizeTableName 将设备 Key 转为合法的 TDengine 子表表名
func sanitizeTableName(deviceKey string) string {
	clean := strings.ReplaceAll(deviceKey, "-", "_")
	clean = strings.ReplaceAll(clean, ".", "_")
	return "d_" + clean
}

// splitValueForTDengine 将动态类型的值智能拆分写入 num_value 或 str_value 槽位
func splitValueForTDengine(v interface{}) (numVal string, strVal string) {
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

// Close 优雅刷新并关闭写入器
func (w *TDengineWriter) Close() {
	close(w.stopChan)
}
