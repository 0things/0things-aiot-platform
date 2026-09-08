package model

import "time"

// DeviceShadow represents the real-time state snapshot of a device.
type DeviceShadow struct {
	DeviceKey  string                 `json:"device_key"`
	Attributes map[string]interface{} `json:"attributes"`
	LastSeen   time.Time              `json:"last_seen"`
}
