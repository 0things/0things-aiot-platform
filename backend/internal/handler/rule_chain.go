package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type RuleChainHandler struct {
	*Handler
	svc service.RuleChainServiceInterface
}

func NewRuleChainHandler(h *Handler, svc service.RuleChainServiceInterface) *RuleChainHandler {
	return &RuleChainHandler{Handler: h, svc: svc}
}

func ruleChainJSON(chain model.RuleChain) v1.RuleChain {
	graph := map[string]any{}
	_ = json.Unmarshal(chain.Graph, &graph)
	return v1.RuleChain{UUID: chain.UUID, OrganizationID: chain.OrganizationID, Name: chain.Name, Description: chain.Description, Status: chain.Status, Version: chain.Version, Graph: graph, CreatedAt: chain.CreatedAt.Format(time.RFC3339), UpdatedAt: chain.UpdatedAt.Format(time.RFC3339)}
}

func bindRuleChainRequest(c *gin.Context) (v1.RuleChainRequest, error) {
	var req v1.RuleChainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, err
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return req, errors.New("rule chain name is required")
	}
	return req, nil
}

// ListRuleChains godoc
// @Summary List rule chains
// @Tags Rule chains
// @Produce json
// @Security Bearer
// @Param request query v1.ListRuleChainsRequest false "Query parameters"
// @Success 200 {object} v1.ApiResponse[v1.ListRuleChainsResponse]
// @Router /rule-chains [get]
func (h *RuleChainHandler) ListRuleChains(c *gin.Context) {
	var req v1.ListRuleChainsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	items, total, err := h.svc.List(c, req.Page, req.PageSize, req.Search, req.Status)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	result := make([]v1.RuleChain, len(items))
	for i := range items {
		result[i] = ruleChainJSON(items[i])
	}
	v1.HandleSuccess(c, v1.ListRuleChainsResponse{Items: result, Total: total, Page: req.Page, PageSize: req.PageSize})
}

// GetRuleChain godoc
// @Summary Get rule chain
// @Tags Rule chains
// @Produce json
// @Security Bearer
// @Param uuid path string true "Rule chain UUID"
// @Success 200 {object} v1.ApiResponse[v1.GetRuleChainResponse]
// @Router /rule-chains/{uuid} [get]
func (h *RuleChainHandler) GetRuleChain(c *gin.Context) {
	chain, err := h.svc.Get(c, c.Param("uuid"))
	if err != nil {
		v1.HandleError(c, http.StatusNotFound, err, nil)
		return
	}
	v1.HandleSuccess(c, v1.GetRuleChainResponse{RuleChain: ruleChainJSON(*chain)})
}

// CreateRuleChain godoc
// @Summary Create rule chain
// @Tags Rule chains
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.RuleChainRequest true "Rule chain"
// @Success 200 {object} v1.ApiResponse[v1.CreateRuleChainResponse]
// @Router /rule-chains [post]
func (h *RuleChainHandler) CreateRuleChain(c *gin.Context) {
	req, err := bindRuleChainRequest(c)
	if err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	graph, _ := json.Marshal(req.Graph)
	chain := &model.RuleChain{Name: req.Name, Description: req.Description, Graph: graph}
	if err := h.svc.Create(c, chain); err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	v1.HandleSuccess(c, v1.CreateRuleChainResponse{RuleChain: ruleChainJSON(*chain)})
}

// UpdateRuleChain godoc
// @Summary Update rule chain
// @Tags Rule chains
// @Accept json
// @Produce json
// @Security Bearer
// @Param uuid path string true "Rule chain UUID"
// @Param request body v1.RuleChainRequest true "Rule chain"
// @Success 200 {object} v1.ApiResponse[v1.UpdateRuleChainResponse]
// @Router /rule-chains/{uuid} [put]
func (h *RuleChainHandler) UpdateRuleChain(c *gin.Context) {
	req, err := bindRuleChainRequest(c)
	if err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	chain, err := h.svc.Get(c, c.Param("uuid"))
	if err != nil {
		v1.HandleError(c, http.StatusNotFound, err, nil)
		return
	}
	graph, _ := json.Marshal(req.Graph)
	chain.Name, chain.Description, chain.Graph = req.Name, req.Description, graph
	if err := h.svc.Update(c, chain); err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	v1.HandleSuccess(c, v1.UpdateRuleChainResponse{RuleChain: ruleChainJSON(*chain)})
}

// DeleteRuleChain godoc
// @Summary Delete rule chain
// @Tags Rule chains
// @Produce json
// @Security Bearer
// @Param uuid path string true "Rule chain UUID"
// @Success 200 {object} v1.ApiResponse[v1.RuleChainSuccessResponse]
// @Router /rule-chains/{uuid} [delete]
func (h *RuleChainHandler) DeleteRuleChain(c *gin.Context) {
	if err := h.svc.Delete(c, c.Param("uuid")); err != nil {
		v1.HandleError(c, http.StatusNotFound, err, nil)
		return
	}
	v1.HandleSuccess(c, v1.RuleChainSuccessResponse{Success: true})
}
