// Package alertv1 owns alert HTTP request and response contracts.
package alertv1

import "time"

type AckAlertRequest struct {
	AckBy string `json:"ackBy"`
}//@name AlertAckAlertRequest

type ResolveAlertRequest struct {
	ResolvedBy string `json:"resolvedBy"`
}//@name AlertResolveAlertRequest

type SuccessResponse struct {
	Success bool `json:"success"`
}//@name AlertSuccessResponse

type Alert struct {
	ID           int64      `json:"id"`
	RuleID       int64      `json:"ruleId"`
	RuleName     string     `json:"ruleName"`
	DeviceKey    string     `json:"deviceKey"`
	Severity     string     `json:"severity"`
	Status       string     `json:"status"`
	Summary      string     `json:"summary"`
	Payload      string     `json:"payload"`
	Count        int        `json:"count"`
	RaisedAt     time.Time  `json:"raisedAt"`
	LastRaisedAt time.Time  `json:"lastRaisedAt"`
	AckAt        *time.Time `json:"ackAt,omitempty"`
	AckBy        string     `json:"ackBy"`
	ResolvedAt   *time.Time `json:"resolvedAt,omitempty"`
	ResolvedBy   string     `json:"resolvedBy"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}//@name Alert

type ListAlertsResponse struct {
	Alerts   []Alert `json:"alerts"`
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}//@name AlertListAlertsResponse

type GetAlertResponse struct {
	Alert Alert `json:"alert"`
}//@name AlertGetAlertResponse
