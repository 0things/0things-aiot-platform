package service

import (
	"context"
	"encoding/json"
	"time"

	"0things/pkg/event"
	"transport-mqtt/internal/adaptor"
	"transport-mqtt/pkg/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type TelemetryService struct {
	adaptor       *adaptor.JsonMqttAdaptor
	eventProducer event.Producer
	logger        *log.Logger
}

func NewTelemetryService(adaptor *adaptor.JsonMqttAdaptor, eventProducer event.Producer, logger *log.Logger) *TelemetryService {
	return &TelemetryService{
		adaptor:       adaptor,
		eventProducer: eventProducer,
		logger:        logger,
	}
}

func (s *TelemetryService) Handle(ctx context.Context, msg mqtt.Message) error {
	dk := extractDeviceKeyFromTopic(msg.Topic())
	if dk == "" {
		s.logger.Warn("could not extract deviceKey from topic", zap.String("topic", msg.Topic()))
		return nil
	}
	pk := extractProductKeyFromTopic(msg.Topic())

	payloads, err := s.adaptor.ConvertToTelemetryPayload(msg.Payload())
	if err != nil || len(payloads) == 0 {
		s.logger.Warn("could not parse telemetry payload via JsonMqttAdaptor",
			zap.String("device_key", dk),
			zap.String("topic", msg.Topic()),
			zap.Error(err),
		)
		return nil
	}

	payloadBytes, err := json.Marshal(payloads)
	if err != nil {
		s.logger.Error("failed to marshal normalized telemetry", zap.Error(err))
		return err
	}

	deviceMsg := event.DeviceMessage{
		DeviceKey:   dk,
		ProductKey:  pk,
		Transport:   event.TransportMQTT,
		MessageType: event.MessageTypeTelemetry,
		Payload:     json.RawMessage(payloadBytes),
		Timestamp:   time.Now().UnixMilli(),
		Headers:     map[string]string{"topic": msg.Topic()},
	}

	s.logger.Info("received and normalized MQTT telemetry message via JsonMqttAdaptor",
		zap.String("topic", msg.Topic()),
		zap.String("device_key", dk),
		zap.Int("records", len(payloads)),
	)

	return s.eventProducer.Publish(ctx, event.TopicDeviceTelemetryReport, &deviceMsg)
}
