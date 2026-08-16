package model

import (
	"encoding/json"
	"time"
)

type SceneLinkageDetail struct {
	ID            int64           `gorm:"primaryKey"`
	SceneID       int64           `gorm:"column:scene_id"`
	TriggerConfig json.RawMessage `gorm:"column:trigger_config;type:json"`
	ActionConfig  json.RawMessage `gorm:"column:action_config;type:json"`
	CreatedAt     time.Time       `gorm:"column:created_at"`
	UpdatedAt     time.Time       `gorm:"column:updated_at"`
}

func (SceneLinkageDetail) TableName() string { return "scene_linkage_detail" }
