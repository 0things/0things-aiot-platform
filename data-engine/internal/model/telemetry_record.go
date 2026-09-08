package model

import "time"

// TelemetryRecord represents a single extracted time-series metric data point.
type TelemetryRecord struct {
	DeviceKey string      `json:"device_key"` // Device key identifier
	Metric    string      `json:"metric"`     // Metric name (e.g., temperature, humidity, voltage)
	Value     interface{} `json:"value"`      // Metric value (float64, int, string, boolean)
	Timestamp time.Time   `json:"timestamp"`  // Data sampling timestamp
}
