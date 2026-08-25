package model

import (
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	Id        int64          `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"not null"`
	Slug      string         `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Organization) TableName() string {
	return "organizations"
}
