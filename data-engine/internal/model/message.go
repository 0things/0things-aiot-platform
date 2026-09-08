package model

import (
	"time"
)

// TelemetryRecord 是规则引擎从上报 Payload 中提取出的单条标准化时序指标。
type TelemetryRecord struct {
	DeviceKey string      `json:"device_key"` // 设备 Key
	Metric    string      `json:"metric"`     // 指标名称（如 temperature, humidity, voltage）
	Value     interface{} `json:"value"`      // 指标数值（float64, int, string, boolean）
	Timestamp time.Time   `json:"timestamp"`  // 采样时间戳
}

// AlarmEvent 是当上报指标违反规则阈值时，规则引擎生成的告警事件对象。
type AlarmEvent struct {
	DeviceKey   string    `json:"device_key"`  // 发生告警的设备 Key
	RuleName    string    `json:"rule_name"`   // 触发的规则名称
	Level       string    `json:"level"`       // 告警级别：CRITICAL / MAJOR / MINOR / WARNING
	Description string    `json:"description"` // 告警详情描述
	Timestamp   time.Time `json:"timestamp"`   // 触发时间
}
