package model

import "time"

// AlarmEvent represents an alert generated when reported metrics violate rule thresholds.
type AlarmEvent struct {
	DeviceKey   string    `json:"device_key"`  // Device key where the alarm originated
	RuleName    string    `json:"rule_name"`   // Rule name triggered
	Level       string    `json:"level"`       // Alarm severity level: CRITICAL / MAJOR / MINOR / WARNING
	Description string    `json:"description"` // Alarm description
	Timestamp   time.Time `json:"timestamp"`   // Trigger timestamp
}
