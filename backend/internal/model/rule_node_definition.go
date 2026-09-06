package model

import "time"

// RuleNodeDefinition describes one system-provided rule node for the visual editor.
type RuleNodeDefinition struct {
	ID            int64     `gorm:"column:id;primaryKey" json:"id"`
	UUID          string    `gorm:"column:uuid;type:varchar(36);not null;uniqueIndex:uk_rule_node_definitions_uuid" json:"uuid"`
	Key           string    `gorm:"column:key;type:varchar(128);not null;uniqueIndex:uk_rule_node_definitions_key" json:"key"`
	Version       int       `gorm:"column:version;not null;default:1" json:"version"`
	Category      string    `gorm:"column:category;type:varchar(32);not null;index:idx_rule_node_definitions_category_sort,priority:1" json:"category"`
	Name          string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Description   string    `gorm:"column:description;type:text" json:"description"`
	Icon          string    `gorm:"column:icon;type:varchar(64)" json:"icon"`
	ConfigSchema  string    `gorm:"column:config_schema;type:text;not null" json:"configSchema"`
	DefaultConfig string    `gorm:"column:default_config;type:text;not null" json:"defaultConfig"`
	InputPorts    string    `gorm:"column:input_ports;type:text;not null" json:"inputPorts"`
	OutputPorts   string    `gorm:"column:output_ports;type:text;not null" json:"outputPorts"`
	ExecutorKey   string    `gorm:"column:executor_key;type:varchar(128);not null" json:"executorKey"`
	Enabled       bool      `gorm:"column:enabled;not null;default:true" json:"enabled"`
	IsSystem      bool      `gorm:"column:is_system;not null;default:false" json:"isSystem"`
	SortOrder     int       `gorm:"column:sort_order;not null;default:0;index:idx_rule_node_definitions_category_sort,priority:2" json:"sortOrder"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (RuleNodeDefinition) TableName() string { return "rule_node_definitions" }
