package v1

// QueryDevicesRequest defines filters for the iot_query_devices MCP tool.
type QueryDevicesRequest struct {
	ProductID *int64 `json:"productId,omitempty"`
	Status    string `json:"status,omitempty"`
	Keyword   string `json:"keyword,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// GetDeviceDetailRequest identifies the target of iot_get_device_detail.
type GetDeviceDetailRequest struct {
	DeviceKey string `json:"deviceKey"`
}

// QueryTelemetryHistoryRequest defines the time range for a telemetry query.
type QueryTelemetryHistoryRequest struct {
	DeviceKey string `json:"deviceKey"`
	Property  string `json:"property"`
	StartTime int64  `json:"startTime,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}
