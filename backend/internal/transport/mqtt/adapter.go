package mqtttransport

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/transport"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// Adapter 将 MQTT 原始消息转换为统一消息，并负责向设备 Topic 下发命令。
type Adapter struct {
	client mqtt.Client
	topic  string
	mu     sync.Mutex
}

func NewAdapter(config *viper.Viper) *Adapter {
	broker := config.GetString("device_gateway.mqtt.broker")
	if broker == "" {
		broker = config.GetString("data.mqtt.broker")
	}
	if broker == "" {
		broker = "tcp://127.0.0.1:1883"
	}
	topic := config.GetString("device_gateway.mqtt.ingress_topic")
	if topic == "" {
		topic = enum.MQTTTopicDevicePropertyPostWildcard
	}
	options := mqtt.NewClientOptions().AddBroker(broker).SetClientID("aiot-gateway-" + uuid.NewString()).SetAutoReconnect(true).SetConnectRetry(true)
	if username := config.GetString("device_gateway.mqtt.username"); username != "" {
		options.SetUsername(username)
	}
	if password := config.GetString("device_gateway.mqtt.password"); password != "" {
		options.SetPassword(password)
	}
	return &Adapter{client: mqtt.NewClient(options), topic: topic}
}

func (a *Adapter) Name() string                      { return "mqtt-device" }
func (a *Adapter) Transport() enum.TransportProtocol { return enum.TransportMQTT }

func (a *Adapter) Start(ctx context.Context, onMessage func(context.Context, transport.DeviceMessage) error) error {
	a.mu.Lock()
	token := a.client.Connect()
	a.mu.Unlock()
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("mqtt connect timeout")
	}
	if err := token.Error(); err != nil {
		return err
	}
	handler := func(_ mqtt.Client, message mqtt.Message) {
		deviceKey := deviceKeyFromTopic(message.Topic())
		if deviceKey == "" {
			return
		}
		if onMessage != nil {
			_ = onMessage(ctx, transport.DeviceMessage{DeviceKey: deviceKey, MessageType: "telemetry", Payload: append([]byte(nil), message.Payload()...), Headers: map[string]string{"topic": message.Topic()}, Timestamp: time.Now().UTC()})
		}
	}
	if token := a.client.Subscribe(a.topic, 1, handler); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	<-ctx.Done()
	return nil
}

func (a *Adapter) Send(ctx context.Context, command transport.Command) error {
	topic := command.Headers["topic"]
	if topic == "" {
		topic = "v1/devices/" + command.DeviceKey + "/commands"
	}
	token := a.client.Publish(topic, 1, false, command.Payload)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		return token.Error()
	}
}

func (a *Adapter) Stop(context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.client.IsConnected() {
		a.client.Disconnect(250)
	}
	return nil
}

func deviceKeyFromTopic(topic string) string {
	parts := strings.Split(topic, "/")
	if len(parts) >= 4 && parts[1] == "sys" {
		return parts[3]
	}
	for i := range parts[:len(parts)-1] {
		if parts[i] == "devices" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}
