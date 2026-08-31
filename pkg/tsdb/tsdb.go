package tsdb

import (
	"context"
	"encoding/json"
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
	Value     interface{} `json:"value"`     // 采样值 (可能为 float64, string, bool, 或 map[string]interface{})
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

// ToTypedValue 将通用 interface{} 智能分流为 4 种基础形态 (*numVal, *strVal, *boolVal, *jsonVal)
func ToTypedValue(v interface{}) (numVal *float64, strVal *string, boolVal *bool, jsonVal *string) {
	if v == nil {
		return nil, nil, nil, nil
	}

	switch val := v.(type) {
	case float64:
		return &val, nil, nil, nil
	case float32:
		f := float64(val)
		return &f, nil, nil, nil
	case int:
		f := float64(val)
		return &f, nil, nil, nil
	case int32:
		f := float64(val)
		return &f, nil, nil, nil
	case int64:
		f := float64(val)
		return &f, nil, nil, nil
	case bool:
		return nil, nil, &val, nil
	case map[string]interface{}, []interface{}:
		// 复合对象/GPS定位/数组/结构体 -> 存入 json_v
		if bytes, err := json.Marshal(val); err == nil {
			s := string(bytes)
			return nil, nil, nil, &s
		}
		s := fmt.Sprintf("%v", val)
		return nil, &s, nil, nil
	case string:
		trimmed := strings.TrimSpace(val)
		// 如果字符串本身是合法的 JSON 对象/数组，优先作为 json_v 处理
		if (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
			(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")) {
			var js interface{}
			if err := json.Unmarshal([]byte(trimmed), &js); err == nil {
				return nil, nil, nil, &trimmed
			}
		}
		return nil, &val, nil, nil
	default:
		// 其他未知复合类型尝试 JSON 序列化
		if bytes, err := json.Marshal(val); err == nil && len(bytes) > 0 && (bytes[0] == '{' || bytes[0] == '[') {
			s := string(bytes)
			return nil, nil, nil, &s
		}
		s := fmt.Sprintf("%v", val)
		return nil, &s, nil, nil
	}
}

// UnmarshalJSONValue 将从数据库取出的 json_v 字符串反序列化为具体的 Go 对象
func UnmarshalJSONValue(raw string) interface{} {
	var result interface{}
	if err := json.Unmarshal([]byte(raw), &result); err == nil {
		return result
	}
	return raw
}
