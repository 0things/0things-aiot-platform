package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"0things/pkg/event"
	"transport-mqtt/pkg/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

// IngressHandler handles incoming MQTT traffic from devices and converts them to NATS events.
type IngressHandler struct {
	eventProducer event.Producer
	logger        *log.Logger
}

func NewIngressHandler(eventProducer event.Producer, logger *log.Logger) *IngressHandler {
	return &IngressHandler{
		eventProducer: eventProducer,
		logger:        logger,
	}
}

// HandleTelemetry parses telemetry payloads and publishes TopicDeviceTelemetryReport.
func (h *IngressHandler) HandleTelemetry(_ mqtt.Client, msg mqtt.Message) {
	deviceKey := ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		h.logger.Warn("could not extract deviceKey from telemetry topic", zap.String("topic", msg.Topic()))
		return
	}
	productKey := ExtractProductKey(msg.Topic())

	deviceMsg := event.DeviceMessage{
		DeviceKey:   deviceKey,
		ProductKey:  productKey,
		Transport:   "mqtt",
		MessageType: "telemetry",
		Payload:     json.RawMessage(msg.Payload()),
		Timestamp:   time.Now().UTC(),
		Headers:     map[string]string{"topic": msg.Topic()},
	}

	h.logger.Info("received MQTT telemetry message",
		zap.String("topic", msg.Topic()),
		zap.String("device_key", deviceKey),
		zap.String("product_key", productKey),
		zap.Int("payload_bytes", len(msg.Payload())),
	)

	if h.eventProducer != nil {
		if err := h.eventProducer.Publish(context.Background(), event.TopicDeviceTelemetryReport, &deviceMsg); err != nil {
			h.logger.Error("failed to publish device telemetry event", zap.Error(err))
		}
	}
}

// HandleDeviceEvent parses device events and publishes TopicDeviceEventReport.
func (h *IngressHandler) HandleDeviceEvent(_ mqtt.Client, msg mqtt.Message) {
	if strings.HasSuffix(msg.Topic(), "/thing/event/property/post") {
		return
	}

	deviceKey := ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		h.logger.Warn("could not extract deviceKey from event topic", zap.String("topic", msg.Topic()))
		return
	}
	productKey := ExtractProductKey(msg.Topic())

	deviceMsg := event.DeviceMessage{
		DeviceKey:   deviceKey,
		ProductKey:  productKey,
		Transport:   "mqtt",
		MessageType: "event",
		Payload:     json.RawMessage(msg.Payload()),
		Timestamp:   time.Now().UTC(),
		Headers:     map[string]string{"topic": msg.Topic()},
	}

	h.logger.Info("received MQTT event message",
		zap.String("topic", msg.Topic()),
		zap.String("device_key", deviceKey),
		zap.String("product_key", productKey),
		zap.Int("payload_bytes", len(msg.Payload())),
	)

	if h.eventProducer != nil {
		if err := h.eventProducer.Publish(context.Background(), event.TopicDeviceEventReport, &deviceMsg); err != nil {
			h.logger.Error("failed to publish device event report", zap.Error(err))
		}
	}
}

// HandleOTAProgress parses device OTA progress reports and publishes TopicOTAProgressReport.
func (h *IngressHandler) HandleOTAProgress(_ mqtt.Client, msg mqtt.Message) {
	deviceKey := ExtractDeviceKey(msg.Topic())
	if deviceKey == "" {
		h.logger.Warn("could not extract deviceKey from OTA topic", zap.String("topic", msg.Topic()))
		return
	}

	var report event.OTAUpgradeReport
	if err := json.Unmarshal(msg.Payload(), &report); err != nil {
		h.logger.Warn("invalid OTA report payload", zap.String("topic", msg.Topic()), zap.Error(err))
		return
	}
	if report.DeviceKey == "" {
		report.DeviceKey = deviceKey
	}
	if report.ProductKey == "" {
		if productKey := ExtractProductKey(msg.Topic()); productKey != "" {
			report.ProductKey = productKey
		}
	}
	if report.ReportedAt.IsZero() {
		report.ReportedAt = time.Now()
	}

	h.logger.Info("received OTA progress report",
		zap.String("topic", msg.Topic()),
		zap.String("device_key", deviceKey),
		zap.Any("report", report),
	)

	if h.eventProducer != nil {
		if err := h.eventProducer.Publish(context.Background(), event.TopicOTAProgressReport, &report); err != nil {
			h.logger.Error("failed to publish OTA progress report event", zap.Error(err))
		}
	}
}

// ExtractDeviceKey extracts deviceKey from topic path.
func ExtractDeviceKey(topic string) string {
	parts := strings.Split(topic, "/")
	if len(parts) >= 4 && parts[1] == "sys" {
		return parts[3]
	}
	if len(parts) >= 6 && parts[1] == "ota" && parts[2] == "device" && (parts[3] == "progress" || parts[3] == "inform") {
		return parts[5]
	}
	return ""
}

// ExtractProductKey extracts productKey from topic path.
func ExtractProductKey(topic string) string {
	parts := strings.Split(topic, "/")
	if len(parts) >= 3 && parts[1] == "sys" {
		return parts[2]
	}
	if len(parts) >= 6 && parts[1] == "ota" && parts[2] == "device" && (parts[3] == "progress" || parts[3] == "inform") {
		return parts[4]
	}
	return ""
}
