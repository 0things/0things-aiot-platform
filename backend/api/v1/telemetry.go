package v1

type TelemetryPoint struct {
	Timestamp int64       `json:"timestamp"`
	Property  string      `json:"property"`
	Value     interface{} `json:"value"`
} //@name TelemetryPoint
