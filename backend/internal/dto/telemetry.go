package dto

import (
	"time"
)

// TelemetryPoint 描述单条时序数据采样点。
type TelemetryPoint struct {
	Timestamp int64       `json:"timestamp"` // 采集毫秒时间戳
	Property  string      `json:"property"`  // 属性标识符 (如 "temperature")
	Value     interface{} `json:"value"`     // 原始值 (数字或字符串)
}

// DeviceShadowSnapshot 描述设备实时的属性快照。
type DeviceShadowSnapshot struct {
	DeviceKey  string                 `json:"device_key"`
	Attributes map[string]interface{} `json:"attributes"`
	LastSeen   time.Time              `json:"last_seen"`
}

// TelemetryQueryReq 时序历史数据查询请求参数。
type TelemetryQueryReq struct {
	DeviceKey string `form:"device_key"`
	Property  string `form:"property" binding:"required"`
	StartTime int64  `form:"start_time"`
	EndTime   int64  `form:"end_time"`
	Limit     int    `form:"limit"`
}
