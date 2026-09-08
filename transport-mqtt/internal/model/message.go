package model

import (
	"encoding/json"
	"time"
)

// DeviceMessage is the normalized MQTT uplink envelope consumed by data-engine.
type DeviceMessage struct {
	DeviceKey   string            `json:"device_key"`
	ProductKey  string            `json:"product_key,omitempty"`
	Transport   string            `json:"transport"`
	MessageType string            `json:"message_type"`
	Payload     json.RawMessage   `json:"payload"`
	Timestamp   time.Time         `json:"timestamp"`
	Headers     map[string]string `json:"headers,omitempty"`
}
