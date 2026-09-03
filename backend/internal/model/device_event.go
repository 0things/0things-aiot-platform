package model

import (
	"time"
)

// DeviceEvent records one event emitted by a device.
type DeviceEvent struct {
	ID              int64     `gorm:"primaryKey"`
	UUID            string    `gorm:"column:uuid;type:varchar(36);uniqueIndex"`
	DeviceID        int64     `gorm:"column:device_id;not null;index:idx_device_events_device_time"`
	EventIdentifier string    `gorm:"column:event_identifier;index:idx_device_events_identifier_time"`
	EventType       string    `gorm:"column:event_type;not null;index:idx_device_events_type_time"`
	EventAt         time.Time `gorm:"column:event_at;not null;index:idx_device_events_device_time;index:idx_device_events_type_time"`
	Data            string    `gorm:"column:data;type:text"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func (DeviceEvent) TableName() string { return "device_events" }
