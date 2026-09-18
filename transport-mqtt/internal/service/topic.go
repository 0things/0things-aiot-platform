package service

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

// extractDeviceKeyFromTopic extracts deviceKey from MQTT topic path (always the last segment).
func extractDeviceKeyFromTopic(value string) string {
	parts := strings.Split(strings.Trim(value, "/"), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return ""
}

// extractProductKeyFromTopic extracts productKey from MQTT topic path (always the second to last segment).
func extractProductKeyFromTopic(value string) string {
	parts := strings.Split(strings.Trim(value, "/"), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return ""
}

// publishDeviceMessage constructs a standard event.DeviceMessage from MQTT message and publishes to the event bus.
func publishDeviceMessage(ctx context.Context, producer event.Producer, logger *log.Logger, msg mqtt.Message, msgType event.MessageType, targetTopic event.Topic, logCategory string) error {
	deviceKey := extractDeviceKeyFromTopic(msg.Topic())
	if deviceKey == "" {
		logger.Warn("could not extract deviceKey from topic", zap.String("topic", msg.Topic()))
		return nil
	}

	deviceMsg := event.DeviceMessage{
		DeviceKey:   deviceKey,
		ProductKey:  extractProductKeyFromTopic(msg.Topic()),
		Transport:   event.TransportMQTT,
		MessageType: msgType,
		Payload:     json.RawMessage(msg.Payload()),
		Timestamp:   time.Now().UnixMilli(),
		Headers:     map[string]string{"topic": msg.Topic()},
	}

	logger.Info("received MQTT "+logCategory+" message",
		zap.String("topic", msg.Topic()),
		zap.String("device_key", deviceKey),
		zap.Int("payload_bytes", len(msg.Payload())),
	)

	if err := producer.Publish(ctx, targetTopic, &deviceMsg); err != nil {
		logger.Error("failed to publish MQTT "+logCategory+" message to event bus",
			zap.String("target_topic", string(targetTopic)),
			zap.String("device_key", deviceKey),
			zap.Error(err),
		)
		return err
	}

	logger.Debug("successfully published MQTT "+logCategory+" message to event bus",
		zap.String("target_topic", string(targetTopic)),
		zap.String("device_key", deviceKey),
	)
	return nil
}
