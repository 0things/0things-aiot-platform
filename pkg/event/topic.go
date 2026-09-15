package event

import (
	"fmt"
	"strings"
)

// Topic represents standard system event topic identifiers.
type Topic string

const (
	// OTA domain events
	TopicOTAUpgradeCommandMQTT Topic = "ota.upgrade.command.mqtt.v1"
	TopicOTAUpgradeCommandHTTP Topic = "ota.upgrade.command.http.v1"
	TopicOTAUpgradeCommand     Topic = "ota.upgrade.command.v1"
	TopicOTAProgressReport     Topic = "ota.progress.report.v1"
	TopicOTADeviceInfo         Topic = "ota.device.info.v1"

	// Device uplink & lifecycle domain events
	TopicDeviceTelemetryReport Topic = "device.telemetry.report.v1"
	TopicDeviceEventReport     Topic = "device.event.report.v1"
	TopicDeviceOnlineStatus    Topic = "device.online.status.v1"
)

// TopicOTAUpgradeCommandByTransport returns the protocol-specific OTA upgrade command topic.
func TopicOTAUpgradeCommandByTransport(transport string) (Topic, error) {
	switch strings.ToLower(strings.TrimSpace(transport)) {
	case "mqtt":
		return TopicOTAUpgradeCommandMQTT, nil
	case "http":
		return TopicOTAUpgradeCommandHTTP, nil
	default:
		return "", fmt.Errorf("unsupported transport for OTA upgrade command: %q", transport)
	}
}

// String returns string representation of Topic.
func (t Topic) String() string {
	return string(t)
}
