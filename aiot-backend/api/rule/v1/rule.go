// Package rulev1 owns the rule HTTP contract and mirrors device-service.
package rulev1

import (
	"encoding/json"
	"time"
)

type RuleRequest struct {
	Name            string          `json:"name" binding:"required"`
	Description     string          `json:"description"`
	Type            string          `json:"type"`
	Status          string          `json:"status"`
	ProductID       int64           `json:"productId"`
	Priority        int             `json:"priority"`
	TriggerConfig   json.RawMessage `json:"triggerConfig"`
	ConditionConfig json.RawMessage `json:"conditionConfig"`
	ActionConfig    json.RawMessage `json:"actionConfig"`
	SQLConfig       json.RawMessage `json:"sqlConfig"`
	Tags            json.RawMessage `json:"tags"`
}//@name RuleRequest

type SuccessResponse struct {
	Success bool `json:"success"`
}//@name RuleSuccessResponse

type Rule struct {
	ID                  int64      `json:"id"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Type                string     `json:"type"`
	Status              string     `json:"status"`
	ProductID           int64      `json:"productId"`
	Priority            int        `json:"priority"`
	TriggerConfig       string     `json:"triggerConfig"`
	ConditionConfig     string     `json:"conditionConfig"`
	ActionConfig        string     `json:"actionConfig"`
	SQLConfig           string     `json:"sqlConfig"`
	ExecutionCount      int64      `json:"executionCount"`
	SuccessCount        int64      `json:"successCount"`
	FailureCount        int64      `json:"failureCount"`
	LastExecutionStatus string     `json:"lastExecutionStatus"`
	CreatedBy           string     `json:"createdBy"`
	Tags                []string   `json:"tags"`
	LastExecutedAt      *time.Time `json:"lastExecutedAt,omitempty"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}//@name Rule

type RuleExecution struct {
	ID              int64     `json:"id"`
	RuleID          int64     `json:"ruleId"`
	RuleName        string    `json:"ruleName"`
	Status          string    `json:"status"`
	TriggeredAt     time.Time `json:"triggeredAt"`
	TriggerData     string    `json:"triggerData"`
	ConditionResult bool      `json:"conditionResult"`
	Duration        int       `json:"duration"`
	Error           string    `json:"error"`
	CreatedAt       time.Time `json:"createdAt"`
}//@name RuleExecution

type AvailableField struct {
	Field       string `json:"field"`
	Type        string `json:"type"`
	Label       string `json:"label"`
	Description string `json:"description"`
}//@name RuleAvailableField

type ListRulesResponse struct {
	Items    []Rule `json:"items"`
	Total    int64  `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}//@name RuleListRulesResponse

type GetRuleResponse struct {
	Rule Rule `json:"rule"`
}//@name RuleGetRuleResponse
type CreateRuleResponse struct {
	Rule Rule `json:"rule"`
}//@name RuleCreateRuleResponse
type UpdateRuleResponse struct {
	Rule Rule `json:"rule"`
}//@name RuleUpdateRuleResponse
type SetRuleStatusResponse struct {
	Rule Rule `json:"rule"`
}//@name RuleSetRuleStatusResponse

type ListRuleExecutionsResponse struct {
	Executions []RuleExecution `json:"executions"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
}//@name RuleListRuleExecutionsResponse

type EvaluateRuleResponse struct {
	Execution RuleExecution `json:"execution"`
}//@name RuleEvaluateRuleResponse
type ListAvailableFieldsResponse struct {
	Fields []AvailableField `json:"fields"`
}//@name RuleListAvailableFieldsResponse
