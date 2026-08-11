package handler

import (
	"encoding/json"
	"errors"
	"strings"

	ruleV1 "0things-backend/api/rule/v1"
	"0things-backend/internal/model"
	"0things-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type RuleHandler struct {
	*Handler
	svc *service.RuleService
}

func NewRuleHandler(h *Handler, svc *service.RuleService) *RuleHandler {
	return &RuleHandler{Handler: h, svc: svc}
}

func ruleJSON(rule model.Rule) ruleV1.Rule {
	var tags []string
	_ = json.Unmarshal(rule.Tags, &tags)
	return ruleV1.Rule{
		ID: rule.ID, Name: rule.Name, Description: rule.Description, Type: rule.Type, Status: rule.Status,
		ProductID: rule.ProductID, Priority: rule.Priority, TriggerConfig: string(rule.TriggerConfig),
		ConditionConfig: string(rule.ConditionConfig), ActionConfig: string(rule.ActionConfig), SQLConfig: string(rule.SQLConfig),
		ExecutionCount: rule.ExecutionCount, SuccessCount: rule.SuccessCount, FailureCount: rule.FailureCount,
		LastExecutionStatus: rule.LastExecutionStatus, CreatedBy: rule.CreatedBy, Tags: tags,
		LastExecutedAt: rule.LastExecutedAt, CreatedAt: rule.CreatedAt, UpdatedAt: rule.UpdatedAt,
	}
}

func ruleExecutionJSON(execution model.RuleExecution) ruleV1.RuleExecution {
	return ruleV1.RuleExecution{
		ID: execution.ID, RuleID: execution.RuleID, RuleName: execution.RuleName, Status: execution.Status,
		TriggeredAt: execution.TriggeredAt, TriggerData: string(execution.TriggerData),
		ConditionResult: execution.ConditionResult, Duration: execution.Duration, Error: execution.Error, CreatedAt: execution.CreatedAt,
	}
}

func (h *RuleHandler) ListRules(c *gin.Context) {
	pageNumber, pageSize := page(c, 20)
	rules, total, err := h.svc.List(c, pageNumber, pageSize, c.Query("type"), c.Query("status"), c.Query("search"))
	if err != nil {
		deviceError(c, err)
		return
	}
	items := make([]ruleV1.Rule, len(rules))
	for i, rule := range rules {
		items[i] = ruleJSON(rule)
	}
	c.JSON(200, ruleV1.ListRulesResponse{Items: items, Total: total, Page: pageNumber, PageSize: pageSize})
}

func (h *RuleHandler) GetRule(c *gin.Context) {
	ruleID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	rule, err := h.svc.Get(c, ruleID)
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, ruleV1.GetRuleResponse{Rule: ruleJSON(*rule)})
}

func (h *RuleHandler) CreateRule(c *gin.Context) {
	var req ruleV1.RuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	rule := &model.Rule{Name: req.Name, Description: req.Description, Type: req.Type, Status: req.Status, ProductID: req.ProductID, Priority: req.Priority, TriggerConfig: req.TriggerConfig, ConditionConfig: req.ConditionConfig, ActionConfig: req.ActionConfig, SQLConfig: req.SQLConfig, Tags: req.Tags}
	if rule.Status == "" {
		rule.Status = "draft"
	}
	if err := h.svc.Create(c, rule); err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, ruleV1.CreateRuleResponse{Rule: ruleJSON(*rule)})
}

func (h *RuleHandler) UpdateRule(c *gin.Context) {
	ruleID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	var req ruleV1.RuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	rule, err := h.svc.Get(c, ruleID)
	if err != nil {
		deviceError(c, err)
		return
	}
	rule.Name, rule.Description, rule.Type, rule.Status = req.Name, req.Description, req.Type, req.Status
	rule.ProductID, rule.Priority = req.ProductID, req.Priority
	rule.TriggerConfig, rule.ConditionConfig, rule.ActionConfig, rule.SQLConfig, rule.Tags = req.TriggerConfig, req.ConditionConfig, req.ActionConfig, req.SQLConfig, req.Tags
	if err := h.svc.Update(c, rule); err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, ruleV1.UpdateRuleResponse{Rule: ruleJSON(*rule)})
}

func (h *RuleHandler) DeleteRule(c *gin.Context) {
	ruleID, err := id(c)
	if err == nil {
		err = h.svc.Delete(c, ruleID)
	}
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, ruleV1.SuccessResponse{Success: true})
}

func (h *RuleHandler) RuleAction(c *gin.Context) {
	rawID := c.Param("id")
	var status string
	var ruleIDText string
	switch {
	case strings.HasSuffix(rawID, ":enable"):
		status = "enabled"
		ruleIDText = strings.TrimSuffix(rawID, ":enable")
	case strings.HasSuffix(rawID, ":disable"):
		status = "disabled"
		ruleIDText = strings.TrimSuffix(rawID, ":disable")
	case strings.HasSuffix(rawID, ":evaluate"):
		h.EvaluateRule(c)
		return
	default:
		deviceError(c, errors.New("unknown rule action"))
		return
	}
	c.Params = gin.Params{{Key: "id", Value: ruleIDText}}
	ruleID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	rule, err := h.svc.SetStatus(c, ruleID, status)
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, ruleV1.SetRuleStatusResponse{Rule: ruleJSON(*rule)})
}

func (h *RuleHandler) RuleExecutions(c *gin.Context) {
	ruleID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	pageNumber, pageSize := page(c, 20)
	executions, total, err := h.svc.ListExecutions(c, ruleID, pageNumber, pageSize)
	if err != nil {
		deviceError(c, err)
		return
	}
	items := make([]ruleV1.RuleExecution, len(executions))
	for i, execution := range executions {
		items[i] = ruleExecutionJSON(execution)
	}
	c.JSON(200, ruleV1.ListRuleExecutionsResponse{Executions: items, Total: total, Page: pageNumber, PageSize: pageSize})
}

func (h *RuleHandler) EvaluateRule(c *gin.Context) {
	ruleID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	execution, err := h.svc.Evaluate(c, ruleID)
	if err != nil {
		deviceError(c, err)
		return
	}
	c.JSON(200, ruleV1.EvaluateRuleResponse{Execution: ruleExecutionJSON(*execution)})
}

func (h *RuleHandler) AvailableFields(c *gin.Context) {
	c.JSON(200, ruleV1.ListAvailableFieldsResponse{Fields: []ruleV1.AvailableField{}})
}
