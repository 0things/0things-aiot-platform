package transport

import "time"

// DeviceMessage 是适配器与设备业务之间交换的统一消息，原始报文由适配器负责解码。
type DeviceMessage struct {
	DeviceID    int64             `json:"deviceId"`
	DeviceKey   string            `json:"deviceKey"`
	EndpointID  int64             `json:"endpointId"`
	MessageType string            `json:"messageType"`
	Payload     []byte            `json:"payload"`
	Headers     map[string]string `json:"headers,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	TraceID     string            `json:"traceId,omitempty"`
}

// Command 是发送到设备端点的统一下行命令，具体协议由适配器选择。
type Command struct {
	DeviceID   int64             `json:"deviceId"`
	DeviceKey  string            `json:"deviceKey"`
	EndpointID int64             `json:"endpointId"`
	Type       string            `json:"type"`
	Payload    []byte            `json:"payload"`
	Headers    map[string]string `json:"headers,omitempty"`
}
