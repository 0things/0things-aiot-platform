package model

import (
	"encoding/json"
	"time"
)

// DeviceMessage 是从 Kafka device.message.v1 消费的标准设备上行消息。
type DeviceMessage struct {
	DeviceKey   string            `json:"device_key"`            // 设备全局唯一标识 Key
	ProductKey  string            `json:"product_key,omitempty"`  // 产品 Key
	Transport   string            `json:"transport"`             // 来源协议：mqtt / http / coap
	MessageType string            `json:"message_type"`          // 消息类型：telemetry / attributes / event / ota_report
	Payload     json.RawMessage   `json:"payload"`               // 原始载荷 JSON
	Timestamp   time.Time         `json:"timestamp"`             // 上报时间（UTC）
	Headers     map[string]string `json:"headers,omitempty"`     // 协议扩展元数据
}

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
