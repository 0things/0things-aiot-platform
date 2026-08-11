package model

import (
	"encoding/json"
	"time"
)

type DeviceShadow struct {
	ID        int64           `gorm:"column:id;primaryKey"`
	DeviceID  int64           `gorm:"column:device_id"`
	Desired   json.RawMessage `gorm:"column:desired;type:json"`
	Reported  json.RawMessage `gorm:"column:reported;type:json"`
	Metadata  json.RawMessage `gorm:"column:metadata;type:json"`
	Version   int64           `gorm:"column:version"`
	CreatedAt time.Time       `gorm:"column:created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
}

func (DeviceShadow) TableName() string { return "device_shadows" }

// DeviceShadowHistory is new.  Unlike the legacy tables it is owned by this
// service and stores an immutable snapshot for the UI history endpoint.
type DeviceShadowHistory struct {
	ID        int64           `gorm:"primaryKey"`
	DeviceID  int64           `gorm:"index;not null"`
	Version   int64           `gorm:"not null"`
	Source    string          `gorm:"size:16;not null"`
	Desired   json.RawMessage `gorm:"type:json"`
	Reported  json.RawMessage `gorm:"type:json"`
	CreatedAt time.Time       `gorm:"index;not null"`
}

func (DeviceShadowHistory) TableName() string { return "device_shadow_histories" }
