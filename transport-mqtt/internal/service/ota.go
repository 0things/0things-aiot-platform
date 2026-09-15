package service

import (
	"context"

	"0things/pkg/event"
	"transport-mqtt/pkg/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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

// HandleProgress processes OTA progress reports (/sys/ota/device/progress/{productKey}/{deviceKey}).
func (s *OTAService) HandleProgress(ctx context.Context, msg mqtt.Message) error {
	return publishDeviceMessage(ctx, s.eventProducer, s.logger, msg, event.MessageTypeOTAProgress, event.TopicOTAProgressReport, "OTA progress")
}

// HandleInform processes OTA device info / version reports (/sys/ota/device/inform/{productKey}/{deviceKey}).
func (s *OTAService) HandleInform(ctx context.Context, msg mqtt.Message) error {
	return publishDeviceMessage(ctx, s.eventProducer, s.logger, msg, event.MessageTypeOTAInform, event.TopicOTADeviceInfo, "OTA inform")
}
