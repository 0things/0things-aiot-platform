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

// ListRules godoc
// @Summary 获取规则列表
// @Schemes
// @Description 分页获取规则列表，支持按 type、status、search 过滤
// @Tags 规则模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param type query string false "规则类型"
// @Param status query string false "状态"
// @Param search query string false "搜索关键字"
// @Success 200 {object} ruleV1.ListRulesResponse
// @Router /rules [get]
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

// GetRule godoc
// @Summary 获取规则详情
// @Schemes
// @Description 通过规则 ID 获取规则详情
// @Tags 规则模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "规则 ID"
// @Success 200 {object} ruleV1.GetRuleResponse
// @Router /rules/{id} [get]
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

// CreateRule godoc
// @Summary 创建规则
// @Schemes
// @Description 创建一条新的规则
// @Tags 规则模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body ruleV1.RuleRequest true "params"
// @Success 200 {object} ruleV1.CreateRuleResponse
// @Router /rules [post]
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

// UpdateRule godoc
// @Summary 更新规则
// @Schemes
// @Description 通过规则 ID 更新规则
// @Tags 规则模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "规则 ID"
// @Param request body ruleV1.RuleRequest true "params"
// @Success 200 {object} ruleV1.UpdateRuleResponse
// @Router /rules/{id} [put]
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

// DeleteRule godoc
// @Summary 删除规则
// @Schemes
// @Description 通过规则 ID 删除规则
// @Tags 规则模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "规则 ID"
// @Success 200 {object} ruleV1.SuccessResponse
// @Router /rules/{id} [delete]
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

// RuleAction godoc
// @Summary 规则操作（启用/禁用/评估）
// @Schemes
// @Description 通过路径参数 id 的后缀执行不同操作：`{id}:enable` 启用、`{id}:disable` 禁用、`{id}:evaluate` 立即评估
// @Tags 规则模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "规则 ID 与动作后缀，例如 1:enable"
// @Success 200 {object} ruleV1.SetRuleStatusResponse
// @Router /rules/{id} [post]
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

// RuleExecutions godoc
// @Summary 获取规则执行记录
// @Schemes
// @Description 分页获取指定规则的执行记录
// @Tags 规则模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "规则 ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} ruleV1.ListRuleExecutionsResponse
// @Router /rules/{id}/executions [get]
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

// EvaluateRule godoc
// @Summary 立即评估规则
// @Schemes
// @Description 通过规则 ID 立即触发一次评估并返回执行结果。由 RuleAction 通过 `:evaluate` 后缀间接调用
// @Tags 规则模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "规则 ID"
// @Success 200 {object} ruleV1.EvaluateRuleResponse
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

// AvailableFields godoc
// @Summary 获取规则可用字段
// @Schemes
// @Description 获取规则引擎支持的可用字段列表
// @Tags 规则模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} ruleV1.ListAvailableFieldsResponse
// @Router /rules/available-fields [get]
func (h *RuleHandler) AvailableFields(c *gin.Context) {
	c.JSON(200, ruleV1.ListAvailableFieldsResponse{Fields: []ruleV1.AvailableField{}})
}
