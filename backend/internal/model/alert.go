package model

import (
	"time"
)

type Alert struct {
	ID           int64  `gorm:"primaryKey"`
	RuleID       int64  `gorm:"column:rule_id"`
	RuleName     string `gorm:"column:rule_name"`
	DeviceKey    string `gorm:"column:device_key"`
	Severity     string
	Status       string
	Summary      string
	Payload      string `gorm:"type:text"`
	Fingerprint  string
	Count        int
	RaisedAt     time.Time  `gorm:"column:raised_at"`
	LastRaisedAt time.Time  `gorm:"column:last_raised_at"`
	AckAt        *time.Time `gorm:"column:ack_at"`
	AckBy        string     `gorm:"column:ack_by"`
	ResolvedAt   *time.Time `gorm:"column:resolved_at"`
	ResolvedBy   string     `gorm:"column:resolved_by"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Alert) TableName() string { return "alerts" }
