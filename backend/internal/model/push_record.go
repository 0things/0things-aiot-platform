package model

import "time"

type DevicePushRecord struct {
	ID            int64  `gorm:"primaryKey"`
	DeviceID      int64  `gorm:"column:device_id"`
	OperationType string `gorm:"column:operation_type"`
	OperationName string `gorm:"column:operation_name"`
	Payload       string
	Status        string
	ErrorMessage  string `gorm:"column:error_message"`
	CreatedBy     string `gorm:"column:created_by"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (DevicePushRecord) TableName() string { return "device_push_records" }
