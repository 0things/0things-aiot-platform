package model

import (
	"time"

	"gorm.io/gorm"
)

type SceneLinkage struct {
	ID          int64           `gorm:"primaryKey"`
	OrganizationID    int64           `gorm:"column:organization_id"`
	Name        string          `gorm:"column:name"`
	Description string          `gorm:"column:description"`
	Enable      int             `gorm:"column:enable"` // 1 = enabled, 0 = disabled
	CreatedAt   time.Time       `gorm:"column:created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"column:deleted_at"`
}

func (SceneLinkage) TableName() string { return "scene_linkage" }
