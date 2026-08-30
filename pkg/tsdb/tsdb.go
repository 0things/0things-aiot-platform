package tsdb

import (
	"context"
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
