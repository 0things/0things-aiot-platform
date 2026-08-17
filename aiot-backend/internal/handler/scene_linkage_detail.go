package handler

import (
	"encoding/json"

	sceneLinkageV1 "0things-backend/api/scene_linkage/v1"
	"0things-backend/internal/model"
	"0things-backend/internal/service"
	v1 "0things-backend/api/v1"
	"github.com/gin-gonic/gin"
)

type SceneLinkageDetailHandler struct {
	*Handler
	svc service.SceneLinkageDetailServiceInterface
}

func NewSceneLinkageDetailHandler(h *Handler, svc service.SceneLinkageDetailServiceInterface) *SceneLinkageDetailHandler {
	return &SceneLinkageDetailHandler{Handler: h, svc: svc}
}

func sceneLinkageDetailJSON(detail model.SceneLinkageDetail) sceneLinkageV1.SceneLinkageDetail {
	return sceneLinkageV1.SceneLinkageDetail{
		SceneID:       detail.SceneID,
		TriggerConfig: detail.TriggerConfig,
		ActionConfig:  detail.ActionConfig,
	}
}

// GetSceneLinkageDetail godoc
// @Summary 获取场景联动详情配置
// @Schemes
// @Description 通过场景联动 ID 获取触发器与动作配置
// @Tags 场景联动模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "场景联动 ID"
// @Success 200 {object} v1.ApiResponse[sceneLinkageV1.GetSceneLinkageDetailResponse]
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
	v1.HandleSuccess(c, sceneLinkageV1.GetSceneLinkageDetailResponse{Detail: sceneLinkageDetailJSON(*detail)})
}

// CreateSceneLinkageDetail godoc
// @Summary 创建场景联动详情配置
// @Schemes
// @Description 为指定场景联动创建触发器与动作配置
// @Tags 场景联动模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "场景联动 ID"
// @Param request body sceneLinkageV1.SceneLinkageDetailRequest true "params"
// @Success 200 {object} v1.ApiResponse[sceneLinkageV1.CreateSceneLinkageDetailResponse]
// @Router /scene-linkages/{id}/detail [post]
func (h *SceneLinkageDetailHandler) CreateSceneLinkageDetail(c *gin.Context) {
	sceneID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	var req sceneLinkageV1.SceneLinkageDetailRequest
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
	v1.HandleSuccess(c, sceneLinkageV1.CreateSceneLinkageDetailResponse{Detail: sceneLinkageDetailJSON(*detail)})
}

// UpdateSceneLinkageDetail godoc
// @Summary 更新场景联动详情配置
// @Schemes
// @Description 通过场景联动 ID 更新触发器与动作配置
// @Tags 场景联动模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "场景联动 ID"
// @Param request body sceneLinkageV1.SceneLinkageDetailRequest true "params"
// @Success 200 {object} v1.ApiResponse[sceneLinkageV1.UpdateSceneLinkageDetailResponse]
// @Router /scene-linkages/{id}/detail [put]
func (h *SceneLinkageDetailHandler) UpdateSceneLinkageDetail(c *gin.Context) {
	sceneID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	var req sceneLinkageV1.SceneLinkageDetailRequest
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
	v1.HandleSuccess(c, sceneLinkageV1.UpdateSceneLinkageDetailResponse{Detail: sceneLinkageDetailJSON(*detail)})
}
