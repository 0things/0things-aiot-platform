package event

// Topic represents standard system event topic identifiers.
type Topic string

const (
	// OTA domain events
	TopicOTAUpgradeCommand Topic = "ota.upgrade.command.v1"
	TopicOTAProgressReport Topic = "ota.progress.report.v1"

	// Device uplink & lifecycle domain events
	TopicDeviceTelemetryReport Topic = "device.telemetry.report.v1"
	TopicDeviceAttributeReport Topic = "device.attribute.report.v1"
	TopicDeviceEventReport     Topic = "device.event.report.v1"
	TopicDeviceOnlineStatus    Topic = "device.online.status.v1"

	// Rule engine & alerting domain events
	TopicAlarmTriggered Topic = "alarm.triggered.v1"
)

func (t Topic) String() string {
	return string(t)
}
