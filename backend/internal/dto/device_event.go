package dto

import (
	"time"
)

// ListDeviceEventsQuery defines internal query parameters for querying device events from repository.
type ListDeviceEventsQuery struct {
	Page      int    // Current page number
	PageSize  int    // Current page size
	Keyword   string // Search keyword
	DeviceKey string // Filter by device key
	EventType string // Filter by event type
	StartAt   *int64 // Filter events starting from timestamp (milliseconds)
	EndAt     *int64 // Filter events up to timestamp (milliseconds)
}

// DeviceEventListItem contains an event and resolved device identity for list views.
type DeviceEventListItem struct {
	ID              int64     `gorm:"column:id"`
	UUID            string    `gorm:"column:uuid"`
	DeviceKey       string    `gorm:"column:device_key"`
	DeviceName      string    `gorm:"column:device_name"`
	ProductKey      string    `gorm:"column:product_key"`
	EventIdentifier string    `gorm:"column:event_identifier"`
	EventType       string    `gorm:"column:event_type"`
	EventAt         int64     `gorm:"column:event_at"`
	Data            string    `gorm:"column:data"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}
