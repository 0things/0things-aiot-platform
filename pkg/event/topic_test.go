package event

import (
	"testing"
)

func TestTopicOTAUpgradeCommandByTransport(t *testing.T) {
	tests := []struct {
		transport string
		expected  Topic
	}{
		{"mqtt", TopicOTAUpgradeCommandMQTT},
		{"MQTT", TopicOTAUpgradeCommandMQTT},
		{"http", TopicOTAUpgradeCommandHTTP},
		{"HTTP", TopicOTAUpgradeCommandHTTP},
		{"", TopicOTAUpgradeCommandMQTT},
		{"ws", Topic("ota.upgrade.command.ws.v1")},
	}

	for _, tt := range tests {
		actual := TopicOTAUpgradeCommandByTransport(tt.transport)
		if actual != tt.expected {
			t.Errorf("TopicOTAUpgradeCommandByTransport(%q) = %q, expected %q", tt.transport, actual, tt.expected)
		}
	}
}

func TestTopic_String(t *testing.T) {
	top := TopicDeviceTelemetryReport
	if top.String() != "device.telemetry.report.v1" {
		t.Errorf("unexpected String() = %s", top.String())
	}
}
