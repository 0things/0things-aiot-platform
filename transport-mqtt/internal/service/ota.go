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

// OTAService processes device OTA progress reports and inform version reports.
type OTAService struct {
	eventProducer event.Producer
	logger        *log.Logger
}

// NewOTAService creates an OTAService instance.
func NewOTAService(eventProducer event.Producer, logger *log.Logger) *OTAService {
	return &OTAService{eventProducer: eventProducer, logger: logger}
}

// HandleProgress processes OTA progress reports (/ota/device/progress/+/+ or /sys/+/+/ota/device/progress).
func (s *OTAService) HandleProgress(ctx context.Context, msg mqtt.Message) error {
	deviceKey := topic.ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from OTA progress topic", zap.String("topic", msg.Topic()))
		return nil
	}

	deviceMsg := event.DeviceMessage{
		DeviceKey:   deviceKey,
		ProductKey:  topic.ExtractProductKey(msg.Topic()),
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeOTAProgress,
		Payload:     json.RawMessage(msg.Payload()),
		Timestamp:   time.Now().UnixMilli(),
		Headers:     map[string]string{"topic": msg.Topic()},
	}

	s.logger.Info("received MQTT OTA progress message",
		zap.String("topic", msg.Topic()),
		zap.String("device_key", deviceKey),
		zap.Int("payload_bytes", len(msg.Payload())),
	)
	if s.eventProducer == nil {
		return nil
	}
	return s.eventProducer.Publish(ctx, event.TopicOTAProgressReport, &deviceMsg)
}

// HandleInform processes OTA device info / version reports (/ota/device/inform/+/+).
func (s *OTAService) HandleInform(ctx context.Context, msg mqtt.Message) error {
	deviceKey := topic.ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from OTA inform topic", zap.String("topic", msg.Topic()))
		return nil
	}

	deviceMsg := event.DeviceMessage{
		DeviceKey:   deviceKey,
		ProductKey:  topic.ExtractProductKey(msg.Topic()),
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeOTAInform,
		Payload:     json.RawMessage(msg.Payload()),
		Timestamp:   time.Now().UnixMilli(),
		Headers:     map[string]string{"topic": msg.Topic()},
	}

	s.logger.Info("received MQTT OTA inform message",
		zap.String("topic", msg.Topic()),
		zap.String("device_key", deviceKey),
		zap.Int("payload_bytes", len(msg.Payload())),
	)
	if s.eventProducer == nil {
		return nil
	}
	return s.eventProducer.Publish(ctx, event.TopicOTADeviceInfo, &deviceMsg)
}
