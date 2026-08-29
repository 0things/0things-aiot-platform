package protocolv1

import "time"

type DeviceEndpointRequest struct {
	ProductProtocolID int64  `json:"productProtocolId" binding:"required"`
	Endpoint          string `json:"endpoint"`
	Credentials       string `json:"credentials"`
	ProtocolConfig    string `json:"protocolConfig"`
	Enabled           bool   `json:"enabled"`
}

type DeviceEndpoint struct {
	ID                int64      `json:"id"`
	DeviceID          int64      `json:"deviceId"`
	ProductProtocolID int64      `json:"productProtocolId"`
	Endpoint          string     `json:"endpoint"`
	ProtocolConfig    string     `json:"protocolConfig"`
	Enabled           bool       `json:"enabled"`
	Status            string     `json:"status"`
	LastSeenAt        *time.Time `json:"lastSeenAt"`
	LastError         string     `json:"lastError"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

// DeviceEndpoints 是根据设备所属产品协议实时生成的接入端点。
// 各传输协议使用独立结构，前端可以直接展示或复制对应示例。
type DeviceEndpoints struct {
	HTTP *HTTPEndpoint `json:"http,omitempty"`
	MQTT *MQTTEndpoint `json:"mqtt,omitempty"`
	CoAP *CoAPEndpoint `json:"coap,omitempty"`
}

type HTTPEndpoint struct {
	HTTP         string `json:"http"`
	RPCSubscribe string `json:"rpcSubscribe"`
}

type MQTTEndpoint struct {
	Host                     string `json:"host"`
	Port                     string `json:"port"`
	TelemetryTopic           string `json:"telemetryTopic"`
	AttributesTopic          string `json:"attributesTopic"`
	AttributesSubscribeTopic string `json:"attributesSubscribeTopic"`
	RPCSubscribeTopic        string `json:"rpcSubscribeTopic"`
}

type CoAPEndpoint struct {
	CoAP         string             `json:"coap"`
	Docker       *CoAPDockerExample `json:"docker,omitempty"`
	RPCSubscribe string             `json:"rpcSubscribe"`
}

type CoAPDockerExample struct {
	CoAP string `json:"coap"`
}
