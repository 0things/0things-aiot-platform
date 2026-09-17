package v1

type RuleChainRequest struct {
	Name        string         `json:"name" binding:"required,max=128"`
	Description string         `json:"description"`
	Graph       map[string]any `json:"graph" binding:"required"`
} //@name RuleChainRequest

type ListRuleChainsRequest struct {
	PageRequest
	Search string `form:"search"`
	Status string `form:"status"`
} //@name ListRuleChainsRequest

type RuleChain struct {
	UUID           string         `json:"uuid"`
	OrganizationID string         `json:"organizationId"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Status         string         `json:"status"`
	Version        int            `json:"version"`
	Graph          map[string]any `json:"graph"`
	CreatedAt      string         `json:"createdAt"`
	UpdatedAt      string         `json:"updatedAt"`
} //@name RuleChain

type ListRuleChainsResponse struct {
	Items    []RuleChain `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
} //@name ListRuleChainsResponse

type GetRuleChainResponse struct {
	RuleChain RuleChain `json:"ruleChain"`
} //@name GetRuleChainResponse
type CreateRuleChainResponse struct {
	RuleChain RuleChain `json:"ruleChain"`
} //@name CreateRuleChainResponse
type UpdateRuleChainResponse struct {
	RuleChain RuleChain `json:"ruleChain"`
} //@name UpdateRuleChainResponse
type RuleChainSuccessResponse struct {
	Success bool `json:"success"`
} //@name RuleChainSuccessResponse
