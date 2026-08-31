package tsdb

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Record 描述写入 TSDB 的单条时序属性数据。
type Record struct {
	DeviceKey string      `json:"device_key"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Timestamp time.Time   `json:"timestamp"`
}

// Point 描述从 TSDB 查询返回的单条历史时序点。
type Point struct {
	Timestamp int64       `json:"timestamp"` // 毫秒时间戳
	Metric    string      `json:"metric"`    // 属性标识符
	Value     interface{} `json:"value"`     // 采样值
}

// QueryFilter 描述时序历史曲线查询参数。
type QueryFilter struct {
	DeviceKey string
	Metric    string
	StartTime int64 // 毫秒起始时间
	EndTime   int64 // 毫秒结束时间
	Limit     int   // 点数上限
}

// Client 定义统一时序数据库客户端接口（读写双向标准，可插拔）。
type Client interface {
	// WriteBatch 批量写入时序点数据
	WriteBatch(ctx context.Context, records []Record) error
	// QueryPoints 区间范围查询时序历史数据点
	QueryPoints(ctx context.Context, filter QueryFilter) ([]Point, error)
	// Close 优雅释放连接池
	Close() error
}

// SanitizeTableName 将设备唯一 Key 转换为合法的数据库表名
func SanitizeTableName(deviceKey string) string {
	clean := strings.ReplaceAll(deviceKey, "-", "_")
	clean = strings.ReplaceAll(clean, ".", "_")
	return "d_" + clean
}

// SplitValue 将任意类型的变量拆分为合法的 SQL 数值槽或字符串槽 (用于 SQL 模板拼接)
func SplitValue(v interface{}) (numVal string, strVal string) {
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

// ToTypedValue 将通用 interface{} 统一解析为强类型指针 (*float64, *string)，供各大官方原生 SDK 批量绑定复用
func ToTypedValue(v interface{}) (numVal *float64, strVal *string) {
	switch val := v.(type) {
	case float64:
		return &val, nil
	case int:
		f := float64(val)
		return &f, nil
	case int64:
		f := float64(val)
		return &f, nil
	case string:
		return nil, &val
	case bool:
		var f float64
		if val {
			f = 1
		}
		return &f, nil
	default:
		s := fmt.Sprintf("%v", val)
		return nil, &s
	}
}
