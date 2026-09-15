package event

import (
	"testing"
)

func TestTopicOTAUpgradeCommandByTransport(t *testing.T) {
	tests := []struct {
		transport string
		expected  Topic
		expectErr bool
	}{
		{"mqtt", TopicOTAUpgradeCommandMQTT, false},
		{"MQTT", TopicOTAUpgradeCommandMQTT, false},
		{"http", TopicOTAUpgradeCommandHTTP, false},
		{"HTTP", TopicOTAUpgradeCommandHTTP, false},
		{"", "", true},
		{"ws", "", true},
		{"coap", "", true},
	}

	for _, tt := range tests {
		actual, err := TopicOTAUpgradeCommandByTransport(tt.transport)
		if tt.expectErr {
			if err == nil {
				t.Errorf("TopicOTAUpgradeCommandByTransport(%q) expected error, got nil", tt.transport)
			}
		} else {
			if err != nil {
				t.Errorf("TopicOTAUpgradeCommandByTransport(%q) unexpected error: %v", tt.transport, err)
			}
			if actual != tt.expected {
				t.Errorf("TopicOTAUpgradeCommandByTransport(%q) = %q, expected %q", tt.transport, actual, tt.expected)
			}
		}
	}
}

func TestTopic_String(t *testing.T) {
	top := TopicDeviceTelemetryReport
	if top.String() != "device.telemetry.report.v1" {
		t.Errorf("unexpected String() = %s", top.String())
	}
}
