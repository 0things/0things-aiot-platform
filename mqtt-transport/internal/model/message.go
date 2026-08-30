package model

import (
	"encoding/json"
	"time"
)

// DeviceMessage 是所有协议上行消息（MQTT/HTTP/CoAP）投递到 Kafka device.message.v1 的统一契约结构。
// 屏蔽不同物理接入协议的差异，供下游 rule-engine 和业务后台统一消费与解析。
type DeviceMessage struct {
	DeviceKey   string            `json:"device_key"`            // 设备全局唯一标识 Key
	ProductKey  string            `json:"product_key,omitempty"`  // 产品 Key（用于物模型版本与属性匹配）
	Transport   string            `json:"transport"`             // 接入物理协议来源：mqtt / http / coap
	MessageType string            `json:"message_type"`          // 上行消息类型：telemetry(遥测) / attributes(属性) / event(事件) / ota_report(OTA进度)
	Payload     json.RawMessage   `json:"payload"`               // 原始上报载荷 JSON（保留原始数据避免二次序列化开销）
	Timestamp   time.Time         `json:"timestamp"`             // 设备上报或网关接收时间（UTC）
	Headers     map[string]string `json:"headers,omitempty"`     // 传输协议元数据（如 Topic、客户端IP、QoS 等）
}

// DeviceCommand 是从 Kafka device.command.v1 消费的云端下行控制指令。
// 由 Web 管理后台或 rule-engine 决策生成，通知具体传输网关推送到物理设备。
type DeviceCommand struct {
	BatchID     string          `json:"batch_id,omitempty"`     // 批次任务 ID（主要用于 OTA 升级批次追踪）
	ProductKey  string          `json:"product_key,omitempty"`  // 所属产品 Key
	DeviceKey   string          `json:"device_key"`            // 目标设备 Key
	Transport   string          `json:"transport"`             // 指定的传输协议（为空或与当前 Transport 匹配时执行下发）
	CommandType string          `json:"command_type"`          // 命令动作类型：ota_upgrade / property_set / rpc
	Endpoint    string          `json:"endpoint,omitempty"`     // 自定义下发 Topic/路由（可选，为空时由各 Transport 自动根据规范生成）
	Payload     json.RawMessage `json:"payload"`               // 发送给设备的实际控制载荷
	Timestamp   time.Time       `json:"timestamp"`             // 命令生成时间
}
