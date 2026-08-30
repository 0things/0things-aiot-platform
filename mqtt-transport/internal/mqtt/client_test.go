package mqtt

import (
	"fmt"
	"testing"

	"mqtt-transport/internal/enum"
)

func TestExtractDeviceKey(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		expected string
	}{
		{
			name:     "standard telemetry topic",
			topic:    "/sys/device_abc_123/thing/event/property/post",
			expected: "device_abc_123",
		},
		{
			name:     "standard ota progress topic",
			topic:    "/sys/sensor_999/ota/device/progress",
			expected: "sensor_999",
		},
		{
			name:     "custom event topic",
			topic:    "/sys/meter_001/thing/event/high_voltage/post",
			expected: "meter_001",
		},
		{
			name:     "invalid topic prefix",
			topic:    "/other/device_abc_123/property/post",
			expected: "",
		},
		{
			name:     "empty topic",
			topic:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ExtractDeviceKey(tt.topic)
			if actual != tt.expected {
				t.Errorf("ExtractDeviceKey(%q) = %q; want %q", tt.topic, actual, tt.expected)
			}
		})
	}
}

func TestDownlinkTopicFormatting(t *testing.T) {
	deviceKey := "test_dev_01"

	propertyTopic := fmt.Sprintf(enum.MQTTTplPropertySet, deviceKey)
	if propertyTopic != "/sys/test_dev_01/thing/service/property/set" {
		t.Errorf("unexpected property downlink topic: %s", propertyTopic)
	}

	otaTopic := fmt.Sprintf(enum.MQTTTplOTAUpgrade, deviceKey)
	if otaTopic != "/sys/test_dev_01/ota/device/upgrade" {
		t.Errorf("unexpected ota downlink topic: %s", otaTopic)
	}
}
