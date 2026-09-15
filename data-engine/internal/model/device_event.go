package model

import "time"

// DeviceEvent represents a thing model business or alarm event record.
type DeviceEvent struct {
	ID              int64     `gorm:"primaryKey;column:id"`
	UUID            string    `gorm:"column:uuid;type:varchar(36);uniqueIndex"`
	DeviceKey       string    `gorm:"column:device_key;type:varchar(64);not null;index:idx_device_events_key_time,priority:1"`
	ProductKey      string    `gorm:"column:product_key;type:varchar(64);index"`
	EventIdentifier string    `gorm:"column:event_identifier;type:varchar(128);index"`
	EventType       string    `gorm:"column:event_type;type:varchar(32);not null;index"`
	EventAt         int64     `gorm:"column:event_at;not null;index:idx_device_events_key_time,priority:2"`
	Data            string    `gorm:"column:data;type:text"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func (DeviceEvent) TableName() string { return "device_events" }
