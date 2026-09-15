package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopicKeyExtraction(t *testing.T) {
	tests := []struct {
		name               string
		topic              string
		expectedProductKey string
		expectedDeviceKey  string
	}{
		{
			name:               "property post topic",
			topic:              "/sys/thing/property/post/prod_sensor_01/dev_temp_1001",
			expectedProductKey: "prod_sensor_01",
			expectedDeviceKey:  "dev_temp_1001",
		},
		{
			name:               "event post topic",
			topic:              "/sys/thing/event/post/prod_camera_02/dev_cam_2002",
			expectedProductKey: "prod_camera_02",
			expectedDeviceKey:  "dev_cam_2002",
		},
		{
			name:               "ota progress topic",
			topic:              "/sys/ota/device/progress/prod_gateway_03/dev_gw_3003",
			expectedProductKey: "prod_gateway_03",
			expectedDeviceKey:  "dev_gw_3003",
		},
		{
			name:               "ota inform topic",
			topic:              "/sys/ota/device/inform/prod_gateway_03/dev_gw_3003",
			expectedProductKey: "prod_gateway_03",
			expectedDeviceKey:  "dev_gw_3003",
		},
		{
			name:               "topic with trailing slash",
			topic:              "/sys/thing/property/post/prod_01/dev_01/",
			expectedProductKey: "prod_01",
			expectedDeviceKey:  "dev_01",
		},
		{
			name:               "topic without leading slash",
			topic:              "sys/thing/property/post/prod_01/dev_01",
			expectedProductKey: "prod_01",
			expectedDeviceKey:  "dev_01",
		},
		{
			name:               "too short topic",
			topic:              "/invalid",
			expectedProductKey: "",
			expectedDeviceKey:  "",
		},
		{
			name:               "empty topic",
			topic:              "",
			expectedProductKey: "",
			expectedDeviceKey:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedDeviceKey, extractDeviceKeyFromTopic(tt.topic))
			assert.Equal(t, tt.expectedProductKey, extractProductKeyFromTopic(tt.topic))
		})
	}
}
