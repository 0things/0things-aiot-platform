package model

import (
	"encoding/json"
	"time"
)

// DeviceMessage 是 CoAP 协议设备上行数据投递到 Kafka device.message.v1 的统一格式。
type DeviceMessage struct {
	DeviceKey   string            `json:"device_key"`            // 设备全局唯一标识 Key
	ProductKey  string            `json:"product_key,omitempty"`  // 产品 Key（可选）
	Transport   string            `json:"transport"`             // 接入物理协议，固定为 "coap"
	MessageType string            `json:"message_type"`          // 上行类型：telemetry / attributes
	Payload     json.RawMessage   `json:"payload"`               // CoAP 报文体 JSON
	Timestamp   time.Time         `json:"timestamp"`             // 接收时间戳（UTC）
	Headers     map[string]string `json:"headers,omitempty"`     // 远端 UDP IP 与端口元数据
}
