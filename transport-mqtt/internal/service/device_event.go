package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"0things/pkg/event"
	"transport-mqtt/internal/topic"
	"transport-mqtt/pkg/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type DeviceEventService struct {
	eventProducer event.Producer
	logger        *log.Logger
}

func NewDeviceEventService(eventProducer event.Producer, logger *log.Logger) *DeviceEventService {
	return &DeviceEventService{eventProducer: eventProducer, logger: logger}
}

func (s *DeviceEventService) Handle(ctx context.Context, msg mqtt.Message) error {
	if strings.HasSuffix(msg.Topic(), "/thing/event/property/post") {
		return nil
	}
	deviceKey := topic.ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from event topic", zap.String("topic", msg.Topic()))
		return nil
	}
	deviceMsg := event.DeviceMessage{DeviceKey: deviceKey, ProductKey: topic.ExtractProductKey(msg.Topic()), Transport: "mqtt", MessageType: "event", Payload: json.RawMessage(msg.Payload()), Timestamp: time.Now().UTC(), Headers: map[string]string{"topic": msg.Topic()}}
	s.logger.Info("received MQTT event message", zap.String("topic", msg.Topic()), zap.String("device_key", deviceKey), zap.Int("payload_bytes", len(msg.Payload())))
	if s.eventProducer == nil {
		return nil
	}
	return s.eventProducer.Publish(ctx, event.TopicDeviceEventReport, &deviceMsg)
}
