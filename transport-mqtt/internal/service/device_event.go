package service

import (
	"context"

	"transport-mqtt/internal/event"
	"transport-mqtt/pkg/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DeviceEventService struct {
	eventProducer event.Producer
	logger        *log.Logger
}

func NewDeviceEventService(eventProducer event.Producer, logger *log.Logger) *DeviceEventService {
	return &DeviceEventService{eventProducer: eventProducer, logger: logger}
}

func (s *DeviceEventService) Handle(ctx context.Context, msg mqtt.Message) error {
	return publishDeviceMessage(ctx, s.eventProducer, s.logger, msg, event.MessageTypeEvent, event.TopicDeviceEventReport, "event")
}
