package handler

import "transport-mqtt/internal/topic"

// ExtractDeviceKey extracts deviceKey from topic path.
func ExtractDeviceKey(value string) string {
	return topic.ExtractDeviceKey(value)
}

// ExtractProductKey extracts productKey from topic path.
func ExtractProductKey(value string) string {
	return topic.ExtractProductKey(value)
}
