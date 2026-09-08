package event

import (
	"encoding/json"
	"time"
)

// DeviceMessage is the standard uplink envelope for device telemetry, attributes, and events.
type DeviceMessage struct {
	DeviceKey   string            `json:"device_key"`
	ProductKey  string            `json:"product_key,omitempty"`
	Transport   string            `json:"transport"`
	MessageType string            `json:"message_type"` // telemetry, attributes, event
	Payload     json.RawMessage   `json:"payload"`
	Timestamp   time.Time         `json:"timestamp"`
	Headers     map[string]string `json:"headers,omitempty"`
}
