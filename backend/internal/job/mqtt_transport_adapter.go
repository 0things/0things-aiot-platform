package job

import (
	"context"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/service"
	"aiot-backend/internal/transport"
)

// mqttTransportAdapter 将现有 MQTTService 接入统一 Transport Adapter 边界。
type mqttTransportAdapter struct{ client service.MQTTServiceInterface }

func newMQTTTransportAdapter(client service.MQTTServiceInterface) transport.Adapter {
	return &mqttTransportAdapter{client: client}
}

// NewMQTTTransportAdapterForWire 为 Wire 提供 MQTT 适配器，避免 transport 包依赖依赖注入框架。
func NewMQTTTransportAdapterForWire(client service.MQTTServiceInterface) transport.Adapter {
	return newMQTTTransportAdapter(client)
}

func (a *mqttTransportAdapter) Name() string                      { return "mqtt-service" }
func (a *mqttTransportAdapter) Transport() enum.TransportProtocol { return enum.TransportMQTT }

// Start 不创建新的连接；共享 MQTTService 已在服务启动阶段完成连接和重连。
func (a *mqttTransportAdapter) Start(context.Context, func(context.Context, transport.DeviceMessage) error) error {
	return nil
}
func (a *mqttTransportAdapter) Stop(context.Context) error { return nil }
func (a *mqttTransportAdapter) Send(ctx context.Context, command transport.Command) error {
	topic := command.Headers["topic"]
	if topic == "" {
		return service.ErrMQTTNotConnected
	}
	return a.client.Publish(ctx, topic, 1, command.Payload)
}
