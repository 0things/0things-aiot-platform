package model

import (
	"encoding/json"
	"time"
)

// DeviceEvent records one event emitted by a device. Device identity (key/name)
// is derived through DeviceID via a join and is intentionally not stored here.
type DeviceEvent struct {
	ID        int64           `gorm:"primaryKey"`
	DeviceID  int64           `gorm:"column:device_id;not null;index:idx_device_events_device_time"`
	EventType string          `gorm:"column:event_type;not null;index:idx_device_events_type_time"`
	EventAt   time.Time       `gorm:"column:event_at;not null;index:idx_device_events_device_time;index:idx_device_events_type_time"`
	Data      json.RawMessage `gorm:"column:data;type:json"`
	CreatedAt time.Time       `gorm:"column:created_at"`

	DeviceKey  string `gorm:"<-:false;column:device_key"`
	DeviceName string `gorm:"<-:false;column:device_name"`
}

func (DeviceEvent) TableName() string { return "device_events" }
