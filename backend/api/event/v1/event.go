package eventv1

import (
	"time"

	apiV1 "aiot-backend/api/v1"
)

// ListDeviceEventsRequest contains query filters for listing device events.
type ListDeviceEventsRequest struct {
	apiV1.PageRequest
	Keyword   string     `form:"keyword"`                                   // Search keyword across deviceKey, deviceName, or eventIdentifier
	DeviceKey string     `form:"deviceKey"`                                 // Device key filter
	EventType string     `form:"eventType"`                                 // Event type filter (e.g. INFO, WARN, ERROR)
	StartAt   *time.Time `form:"startAt" time_format:"2006-01-02 15:04:05"` // Start time filter (yyyy-MM-dd HH:mm:ss)
	EndAt     *time.Time `form:"endAt" time_format:"2006-01-02 15:04:05"`   // End time filter (yyyy-MM-dd HH:mm:ss)
} //@name ListDeviceEventsRequest

// DeviceEvent represents a thing model device event entity.
type DeviceEvent struct {
	ID              int64     `json:"id"`              // Event primary ID
	UUID            string    `json:"uuid"`            // Globally unique event UUID
	DeviceKey       string    `json:"deviceKey"`       // Unique device key
	DeviceName      string    `json:"deviceName"`      // Human-readable device name
	EventIdentifier string    `json:"eventIdentifier"` // Thing model event identifier
	EventType       string    `json:"eventType"`       // Event type (INFO, WARN, ERROR)
	EventAt         time.Time `json:"eventAt"`         // Event timestamp
	Data            string    `json:"data"`            // JSON string payload of event parameters
} //@name DeviceEvent

// ListDeviceEventsResponse represents the paginated response of device events.
type ListDeviceEventsResponse struct {
	Events   []DeviceEvent `json:"events"`   // List of device events
	Total    int64         `json:"total"`    // Total count of matching events
	Page     int           `json:"page"`     // Current page number
	PageSize int           `json:"pageSize"` // Current page size
} //@name DeviceEventListDeviceEventsResponse
