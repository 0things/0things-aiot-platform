package model

import (
	"encoding/json"
	"time"
)

// DeviceMessage 是 HTTP 设备上行数据投递到 Kafka device.message.v1 的标准结构。
type DeviceMessage struct {
	DeviceKey   string            `json:"device_key"`            // 设备标识 Key
	ProductKey  string            `json:"product_key,omitempty"`  // 产品 Key（可选）
	Transport   string            `json:"transport"`             // 接入物理协议，固定为 "http"
	MessageType string            `json:"message_type"`          // 上行类型：telemetry / attributes / event
	Payload     json.RawMessage   `json:"payload"`               // 原始 HTTP 请求体 JSON
	Timestamp   time.Time         `json:"timestamp"`             // 上报时间戳（UTC）
	Headers     map[string]string `json:"headers,omitempty"`     // 请求元数据（Content-Type、来源 IP 等）
}
