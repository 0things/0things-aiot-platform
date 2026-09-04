package handler

import (
	"encoding/json"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type SceneLinkageDetailHandler struct {
	*Handler
	svc service.SceneLinkageDetailServiceInterface
}

func NewSceneLinkageDetailHandler(h *Handler, svc service.SceneLinkageDetailServiceInterface) *SceneLinkageDetailHandler {
	return &SceneLinkageDetailHandler{Handler: h, svc: svc}
}

func sceneLinkageDetailJSON(detail model.SceneLinkageDetail) v1.SceneLinkageDetail {
	return v1.SceneLinkageDetail{
		SceneID:       detail.SceneID,
		TriggerConfig: detail.TriggerConfig,
		ActionConfig:  detail.ActionConfig,
	}
}

// GetSceneLinkageDetail godoc
// @Summary Get scene linkage configuration
// @Schemes
// @Description Returns scene linkage configuration.
// @Tags Scene linkages
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Scene linkage ID"
// @Success 200 {object} v1.ApiResponse[v1.GetSceneLinkageDetailResponse] "Successful response"
// @Router /scene-linkages/{id}/detail [get]
func (h *SceneLinkageDetailHandler) GetSceneLinkageDetail(c *gin.Context) {
	sceneID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	detail, err := h.svc.GetBySceneID(c, sceneID)
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, v1.GetSceneLinkageDetailResponse{Detail: sceneLinkageDetailJSON(*detail)})
}

// CreateSceneLinkageDetail godoc
// @Summary Create scene linkage configuration
// @Schemes
// @Description Creates scene linkage configuration.
// @Tags Scene linkages
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Scene linkage ID"
// @Param request body v1.SceneLinkageDetailRequest true "params"
// @Success 200 {object} v1.ApiResponse[v1.CreateSceneLinkageDetailResponse] "Successful response"
// @Router /scene-linkages/{id}/detail [post]
func (h *SceneLinkageDetailHandler) CreateSceneLinkageDetail(c *gin.Context) {
	sceneID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	var req v1.SceneLinkageDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	detail := &model.SceneLinkageDetail{
		SceneID:       sceneID,
		TriggerConfig: req.TriggerConfig,
		ActionConfig:  req.ActionConfig,
	}
	if len(detail.TriggerConfig) == 0 {
		detail.TriggerConfig = json.RawMessage("{}")
	}
	if len(detail.ActionConfig) == 0 {
		detail.ActionConfig = json.RawMessage("{}")
	}
	if err := h.svc.Create(c, detail); err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, v1.CreateSceneLinkageDetailResponse{Detail: sceneLinkageDetailJSON(*detail)})
}

// UpdateSceneLinkageDetail godoc
// @Summary Update scene linkage configuration
// @Schemes
// @Description Updates scene linkage configuration.
// @Tags Scene linkages
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Scene linkage ID"
// @Param request body v1.SceneLinkageDetailRequest true "params"
// @Success 200 {object} v1.ApiResponse[v1.UpdateSceneLinkageDetailResponse] "Successful response"
// @Router /scene-linkages/{id}/detail [put]
func (h *SceneLinkageDetailHandler) UpdateSceneLinkageDetail(c *gin.Context) {
	sceneID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	var req v1.SceneLinkageDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	detail, err := h.svc.GetBySceneID(c, sceneID)
	if err != nil {
		deviceError(c, err)
		return
	}
	detail.TriggerConfig = req.TriggerConfig
	detail.ActionConfig = req.ActionConfig
	if len(detail.TriggerConfig) == 0 {
		detail.TriggerConfig = json.RawMessage("{}")
	}
	if len(detail.ActionConfig) == 0 {
		detail.ActionConfig = json.RawMessage("{}")
	}
	if err := h.svc.Update(c, detail); err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, v1.UpdateSceneLinkageDetailResponse{Detail: sceneLinkageDetailJSON(*detail)})
}
