package event

import (
	"testing"
)

func TestTopicValues(t *testing.T) {
	tests := []struct {
		topic    Topic
		expected string
	}{
		{TopicOTAUpgradeCommand, "ota.upgrade.command.v1"},
		{TopicOTAProgressReport, "ota.progress.report.v1"},
		{TopicDeviceTelemetryReport, "device.telemetry.report.v1"},
		{TopicDeviceAttributeReport, "device.attribute.report.v1"},
		{TopicDeviceEventReport, "device.event.report.v1"},
		{TopicDeviceOnlineStatus, "device.online.status.v1"},
		{TopicAlarmTriggered, "alarm.triggered.v1"},
	}

	for _, tt := range tests {
		if tt.topic.String() != tt.expected {
			t.Errorf("expected topic string %q, got %q", tt.expected, tt.topic.String())
		}
	}
}
