package v1

import (
	"encoding/json"
	"time"
)

type SceneLinkageRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Enable      int    `json:"enable"`
} //@name SceneLinkageRequest

type ListSceneLinkagesRequest struct {
	PageRequest
	Search string `form:"search"`
	Enable *int   `form:"enable"`
} //@name ListSceneLinkagesRequest

type SceneLinkage struct {
	ID             int64     `json:"id"`
	OrganizationID int64     `json:"organizationId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Enable         int       `json:"enable"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
} //@name SceneLinkage

type ListSceneLinkagesResponse struct {
	Items    []SceneLinkage `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
} //@name ListSceneLinkagesResponse

type GetSceneLinkageResponse struct {
	SceneLinkage SceneLinkage `json:"sceneLinkage"`
} //@name GetSceneLinkageResponse

type CreateSceneLinkageResponse struct {
	SceneLinkage SceneLinkage `json:"sceneLinkage"`
} //@name CreateSceneLinkageResponse

type UpdateSceneLinkageResponse struct {
	SceneLinkage SceneLinkage `json:"sceneLinkage"`
} //@name UpdateSceneLinkageResponse

type SceneLinkageSuccessResponse struct {
	Success bool `json:"success"`
} //@name SceneLinkageSuccessResponse

type SceneLinkageDetailRequest struct {
	TriggerConfig json.RawMessage `json:"triggerConfig"`
	ActionConfig  json.RawMessage `json:"actionConfig"`
} //@name SceneLinkageDetailRequest

type SceneLinkageDetail struct {
	SceneID       int64           `json:"sceneId"`
	TriggerConfig json.RawMessage `json:"triggerConfig"`
	ActionConfig  json.RawMessage `json:"actionConfig"`
} //@name SceneLinkageDetail

type GetSceneLinkageDetailResponse struct {
	Detail SceneLinkageDetail `json:"detail"`
} //@name GetSceneLinkageDetailResponse

type CreateSceneLinkageDetailResponse struct {
	Detail SceneLinkageDetail `json:"detail"`
} //@name CreateSceneLinkageDetailResponse

type UpdateSceneLinkageDetailResponse struct {
	Detail SceneLinkageDetail `json:"detail"`
} //@name UpdateSceneLinkageDetailResponse
