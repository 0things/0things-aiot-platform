package event

import "strings"

// Topic represents standard system event topic identifiers.
type Topic string

const (
	// OTA domain events
	TopicOTAUpgradeCommandMQTT Topic = "ota.upgrade.command.mqtt.v1"
	TopicOTAUpgradeCommandCoAP Topic = "ota.upgrade.command.coap.v1"
	TopicOTAUpgradeCommandHTTP Topic = "ota.upgrade.command.http.v1"
	TopicOTAUpgradeCommand     Topic = "ota.upgrade.command.v1"
	TopicOTAProgressReport     Topic = "ota.progress.report.v1"

	// Device uplink & lifecycle domain events
	TopicDeviceTelemetryReport Topic = "device.telemetry.report.v1"
	TopicDeviceAttributeReport Topic = "device.attribute.report.v1"
	TopicDeviceEventReport     Topic = "device.event.report.v1"
	TopicDeviceOnlineStatus    Topic = "device.online.status.v1"

	// Rule engine & alerting domain events
	TopicAlarmTriggered Topic = "alarm.triggered.v1"
)

// TopicOTAUpgradeCommandByTransport returns the protocol-specific OTA upgrade command topic.
func TopicOTAUpgradeCommandByTransport(transport string) Topic {
	switch strings.ToLower(strings.TrimSpace(transport)) {
	case "mqtt":
		return TopicOTAUpgradeCommandMQTT
	case "coap":
		return TopicOTAUpgradeCommandCoAP
	case "http":
		return TopicOTAUpgradeCommandHTTP
	default:
		if transport == "" {
			return TopicOTAUpgradeCommandMQTT
		}
		return Topic("ota.upgrade.command." + strings.ToLower(strings.TrimSpace(transport)) + ".v1")
	}
}

func (t Topic) String() string {
	return string(t)
}
