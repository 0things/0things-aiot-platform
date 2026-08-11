package model

import (
	"encoding/json"
	"time"
)

type Rule struct {
	ID                  int64 `gorm:"primaryKey"`
	Name                string
	Description         string
	Type                string
	Status              string
	ProductID           int64 `gorm:"column:product_id"`
	Priority            int
	TriggerConfig       json.RawMessage `gorm:"column:trigger_config;type:json"`
	ConditionConfig     json.RawMessage `gorm:"column:condition_config;type:json"`
	ActionConfig        json.RawMessage `gorm:"column:action_config;type:json"`
	SQLConfig           json.RawMessage `gorm:"column:sql_config;type:json"`
	ExecutionCount      int64           `gorm:"column:execution_count"`
	SuccessCount        int64           `gorm:"column:success_count"`
	FailureCount        int64           `gorm:"column:failure_count"`
	LastExecutionStatus string          `gorm:"column:last_execution_status"`
	CreatedBy           string          `gorm:"column:created_by"`
	Tags                json.RawMessage `gorm:"type:json"`
	LastExecutedAt      *time.Time      `gorm:"column:last_executed_at"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (Rule) TableName() string { return "rules" }

type RuleExecution struct {
	ID              int64  `gorm:"primaryKey"`
	RuleID          int64  `gorm:"column:rule_id"`
	RuleName        string `gorm:"column:rule_name"`
	Status          string
	TriggerData     json.RawMessage `gorm:"column:trigger_data;type:json"`
	ConditionResult bool            `gorm:"column:condition_result"`
	Duration        int
	Error           string
	TriggeredAt     time.Time `gorm:"column:triggered_at"`
	CreatedAt       time.Time
}

func (RuleExecution) TableName() string { return "rule_executions" }
