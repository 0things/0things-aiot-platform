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

type OTAProgressService struct {
	eventProducer event.Producer
	logger        *log.Logger
}

func NewOTAProgressService(eventProducer event.Producer, logger *log.Logger) *OTAProgressService {
	return &OTAProgressService{eventProducer: eventProducer, logger: logger}
}

func (s *OTAProgressService) Handle(ctx context.Context, msg mqtt.Message) error {
	deviceKey := topic.ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		s.logger.Warn("could not extract deviceKey from OTA topic", zap.String("topic", msg.Topic()))
		return nil
	}
	var report event.OTAUpgradeReport
	if err := json.Unmarshal(msg.Payload(), &report); err != nil {
		s.logger.Warn("invalid OTA report payload", zap.String("topic", msg.Topic()), zap.Error(err))
		return nil
	}
	if report.DeviceKey == "" {
		report.DeviceKey = deviceKey
	}
	if report.ProductKey == "" {
		report.ProductKey = topic.ExtractProductKey(msg.Topic())
	}
	if report.ReportedAt.IsZero() {
		report.ReportedAt = time.Now()
	}
	s.logger.Info("received OTA progress report", zap.String("topic", msg.Topic()), zap.String("device_key", deviceKey))
	if s.eventProducer == nil {
		return nil
	}
	return s.eventProducer.Publish(ctx, event.TopicOTAProgressReport, &report)
}
