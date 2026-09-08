package event

import (
	"testing"
)

func TestTopicValues(t *testing.T) {
	tests := []struct {
		topic    Topic
		expected string
	}{
		{TopicOTAUpgradeCommandMQTT, "ota.upgrade.command.mqtt.v1"},
		{TopicOTAUpgradeCommandCoAP, "ota.upgrade.command.coap.v1"},
		{TopicOTAUpgradeCommandHTTP, "ota.upgrade.command.http.v1"},
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

func TestTopicOTAUpgradeCommandByTransport(t *testing.T) {
	tests := []struct {
		transport string
		expected  Topic
	}{
		{"mqtt", TopicOTAUpgradeCommandMQTT},
		{"MQTT", TopicOTAUpgradeCommandMQTT},
		{"coap", TopicOTAUpgradeCommandCoAP},
		{"http", TopicOTAUpgradeCommandHTTP},
		{"", TopicOTAUpgradeCommandMQTT},
		{"websocket", "ota.upgrade.command.websocket.v1"},
	}

	for _, tt := range tests {
		actual := TopicOTAUpgradeCommandByTransport(tt.transport)
		if actual != tt.expected {
			t.Errorf("TopicOTAUpgradeCommandByTransport(%q) = %q; want %q", tt.transport, actual, tt.expected)
		}
	}
}
