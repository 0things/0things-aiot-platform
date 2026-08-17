package handler

import (
	sceneLinkageV1 "0things-backend/api/scene_linkage/v1"
	v1 "0things-backend/api/v1"
	"0things-backend/internal/model"
	"0things-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type SceneLinkageHandler struct {
	*Handler
	svc service.SceneLinkageServiceInterface
}

func NewSceneLinkageHandler(h *Handler, svc service.SceneLinkageServiceInterface) *SceneLinkageHandler {
	return &SceneLinkageHandler{Handler: h, svc: svc}
}

func sceneLinkageJSON(scene model.SceneLinkage) sceneLinkageV1.SceneLinkage {
	return sceneLinkageV1.SceneLinkage{
		ID:          scene.ID,
		TenantID:    scene.TenantID,
		Name:        scene.Name,
		Description: scene.Description,
		Enable:      scene.Enable,
		CreatedAt:   scene.CreatedAt,
		UpdatedAt:   scene.UpdatedAt,
	}
}

func parseEnable(c *gin.Context) int {
	switch c.Query("enable") {
	case "1":
		return 1
	case "0":
		return 0
	default:
		return -1
	}
}

// ListSceneLinkages godoc
// @Summary 获取场景联动列表
// @Schemes
// @Description 分页获取场景联动列表，支持按 search、enable 过滤
// @Tags 场景联动模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param search query string false "搜索关键字"
// @Param enable query int false "启用状态：1 启用，0 停用"
// @Success 200 {object} v1.ApiResponse[sceneLinkageV1.ListSceneLinkagesResponse]
// @Router /scene-linkages [get]
func (h *SceneLinkageHandler) ListSceneLinkages(c *gin.Context) {
	pageNumber, pageSize := page(c, 20)
	scenes, total, err := h.svc.List(c, pageNumber, pageSize, c.Query("search"), parseEnable(c))
	if err != nil {
		deviceError(c, err)
		return
	}
	items := make([]sceneLinkageV1.SceneLinkage, len(scenes))
	for i, scene := range scenes {
		items[i] = sceneLinkageJSON(scene)
	}
	v1.HandleSuccess(c, sceneLinkageV1.ListSceneLinkagesResponse{Items: items, Total: total, Page: pageNumber, PageSize: pageSize})
}

// GetSceneLinkage godoc
// @Summary 获取场景联动详情
// @Schemes
// @Description 通过 ID 获取场景联动
// @Tags 场景联动模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "场景联动 ID"
// @Success 200 {object} v1.ApiResponse[sceneLinkageV1.GetSceneLinkageResponse]
// @Router /scene-linkages/{id} [get]
func (h *SceneLinkageHandler) GetSceneLinkage(c *gin.Context) {
	sceneID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	scene, err := h.svc.Get(c, sceneID)
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, sceneLinkageV1.GetSceneLinkageResponse{SceneLinkage: sceneLinkageJSON(*scene)})
}

// CreateSceneLinkage godoc
// @Summary 创建场景联动
// @Schemes
// @Description 创建一条新的场景联动
// @Tags 场景联动模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body sceneLinkageV1.SceneLinkageRequest true "params"
// @Success 200 {object} v1.ApiResponse[sceneLinkageV1.CreateSceneLinkageResponse]
// @Router /scene-linkages [post]
func (h *SceneLinkageHandler) CreateSceneLinkage(c *gin.Context) {
	var req sceneLinkageV1.SceneLinkageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	scene := &model.SceneLinkage{Name: req.Name, Description: req.Description, Enable: req.Enable}
	if scene.Enable != 0 && scene.Enable != 1 {
		scene.Enable = 1
	}
	if err := h.svc.Create(c, scene); err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, sceneLinkageV1.CreateSceneLinkageResponse{SceneLinkage: sceneLinkageJSON(*scene)})
}

// UpdateSceneLinkage godoc
// @Summary 更新场景联动
// @Schemes
// @Description 通过 ID 更新场景联动
// @Tags 场景联动模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "场景联动 ID"
// @Param request body sceneLinkageV1.SceneLinkageRequest true "params"
// @Success 200 {object} v1.ApiResponse[sceneLinkageV1.UpdateSceneLinkageResponse]
// @Router /scene-linkages/{id} [put]
func (h *SceneLinkageHandler) UpdateSceneLinkage(c *gin.Context) {
	sceneID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	var req sceneLinkageV1.SceneLinkageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	scene, err := h.svc.Get(c, sceneID)
	if err != nil {
		deviceError(c, err)
		return
	}
	scene.Name, scene.Description, scene.Enable = req.Name, req.Description, req.Enable
	if scene.Enable != 0 && scene.Enable != 1 {
		scene.Enable = 1
	}
	if err := h.svc.Update(c, scene); err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, sceneLinkageV1.UpdateSceneLinkageResponse{SceneLinkage: sceneLinkageJSON(*scene)})
}

// DeleteSceneLinkage godoc
// @Summary 删除场景联动
// @Schemes
// @Description 通过 ID 删除场景联动
// @Tags 场景联动模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "场景联动 ID"
// @Success 200 {object} v1.ApiResponse[sceneLinkageV1.SuccessResponse]
// @Router /scene-linkages/{id} [delete]
func (h *SceneLinkageHandler) DeleteSceneLinkage(c *gin.Context) {
	sceneID, err := id(c)
	if err == nil {
		err = h.svc.Delete(c, sceneID)
	}
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, sceneLinkageV1.SuccessResponse{Success: true})
}
