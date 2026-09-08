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
			topic:    "/sys/prod_xyz/device_abc_123/thing/event/property/post",
			expected: "device_abc_123",
		},
		{
			name:     "standard ota progress topic",
			topic:    "/sys/prod_xyz/sensor_999/ota/device/progress",
			expected: "sensor_999",
		},
		{
			name:     "custom event topic",
			topic:    "/sys/prod_xyz/meter_001/thing/event/high_voltage/post",
			expected: "meter_001",
		},
		{
			name:     "independent ota progress topic",
			topic:    "/ota/device/progress/prod_xyz/sensor_999",
			expected: "sensor_999",
		},
		{
			name:     "invalid topic prefix",
			topic:    "/other/prod_xyz/device_abc_123/property/post",
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

func TestExtractProductKey(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		expected string
	}{
		{
			name:     "standard telemetry topic",
			topic:    "/sys/prod_xyz/device_abc_123/thing/event/property/post",
			expected: "prod_xyz",
		},
		{
			name:     "standard ota progress topic",
			topic:    "/sys/prod_xyz/sensor_999/ota/device/progress",
			expected: "prod_xyz",
		},
		{
			name:     "independent ota inform topic",
			topic:    "/ota/device/inform/prod_xyz/sensor_999",
			expected: "prod_xyz",
		},
		{
			name:     "invalid topic prefix",
			topic:    "/other/prod_xyz/device_abc_123/property/post",
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
			actual := ExtractProductKey(tt.topic)
			if actual != tt.expected {
				t.Errorf("ExtractProductKey(%q) = %q; want %q", tt.topic, actual, tt.expected)
			}
		})
	}
}

func TestDownlinkTopicFormatting(t *testing.T) {
	productKey := "prod_xyz"
	deviceKey := "test_dev_01"

	propertyTopic := fmt.Sprintf(enum.MQTTTplPropertySet, productKey, deviceKey)
	if propertyTopic != "/sys/prod_xyz/test_dev_01/thing/service/property/set" {
		t.Errorf("unexpected property downlink topic: %s", propertyTopic)
	}

	otaTopic := fmt.Sprintf(enum.MQTTTplOTAUpgrade, productKey, deviceKey)
	if otaTopic != "/sys/prod_xyz/test_dev_01/ota/device/upgrade" {
		t.Errorf("unexpected ota downlink topic: %s", otaTopic)
	}

	otaTopicV1 := fmt.Sprintf(enum.MQTTTplOTAUpgradeV1, productKey, deviceKey)
	if otaTopicV1 != "/ota/device/upgrade/prod_xyz/test_dev_01" {
		t.Errorf("unexpected ota v1 downlink topic: %s", otaTopicV1)
	}
}
