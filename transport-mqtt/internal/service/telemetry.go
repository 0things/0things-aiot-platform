package service

import (
	"context"
	"encoding/json"
	"time"

	"0things/pkg/event"
	"transport-mqtt/internal/topic"
	"transport-mqtt/pkg/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type TelemetryService struct {
	eventProducer event.Producer
	logger        *log.Logger
}

func NewTelemetryService(eventProducer event.Producer, logger *log.Logger) *TelemetryService {
	return &TelemetryService{eventProducer: eventProducer, logger: logger}
}

func (s *TelemetryService) Handle(ctx context.Context, msg mqtt.Message) error {
	deviceKey := topic.ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from telemetry topic", zap.String("topic", msg.Topic()))
		return nil
	}
	deviceMsg := event.DeviceMessage{
		DeviceKey:   deviceKey,
		ProductKey:  topic.ExtractProductKey(msg.Topic()),
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeTelemetry,
		Payload:     json.RawMessage(msg.Payload()),
		Timestamp:   time.Now().UnixMilli(),
		Headers:     map[string]string{"topic": msg.Topic()},
	}
	s.logger.Info("received MQTT telemetry message", zap.String("topic", msg.Topic()), zap.String("device_key", deviceKey), zap.Int("payload_bytes", len(msg.Payload())))
	if s.eventProducer == nil {
		return nil
	}
	return s.eventProducer.Publish(ctx, event.TopicDeviceTelemetryReport, &deviceMsg)
}
