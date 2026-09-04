package handler

import (
	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type SceneLinkageHandler struct {
	*Handler
	svc service.SceneLinkageServiceInterface
}

func NewSceneLinkageHandler(h *Handler, svc service.SceneLinkageServiceInterface) *SceneLinkageHandler {
	return &SceneLinkageHandler{Handler: h, svc: svc}
}

func sceneLinkageJSON(scene model.SceneLinkage) v1.SceneLinkage {
	return v1.SceneLinkage{
		ID:             scene.ID,
		OrganizationID: scene.OrganizationID,
		Name:           scene.Name,
		Description:    scene.Description,
		Enable:         scene.Enable,
		CreatedAt:      scene.CreatedAt,
		UpdatedAt:      scene.UpdatedAt,
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
// @Summary List scene linkages
// @Schemes
// @Description Lists scene linkages.
// @Tags Scene linkages
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query v1.ListSceneLinkagesRequest false "Query parameters"
// @Success 200 {object} v1.ApiResponse[v1.ListSceneLinkagesResponse] "Successful response"
// @Router /scene-linkages [get]
func (h *SceneLinkageHandler) ListSceneLinkages(c *gin.Context) {
	var req v1.ListSceneLinkagesRequest
	_ = c.ShouldBindQuery(&req)
	pageNumber, pageSize := pageRequest(req.PageRequest, 20)
	enable := -1
	if req.Enable != nil {
		enable = *req.Enable
	}
	scenes, total, err := h.svc.List(c, pageNumber, pageSize, req.Search, enable)
	if err != nil {
		deviceError(c, err)
		return
	}
	items := make([]v1.SceneLinkage, len(scenes))
	for i, scene := range scenes {
		items[i] = sceneLinkageJSON(scene)
	}
	v1.HandleSuccess(c, v1.ListSceneLinkagesResponse{Items: items, Total: total, Page: pageNumber, PageSize: pageSize})
}

// GetSceneLinkage godoc
// @Summary Get scene linkage
// @Schemes
// @Description Returns scene linkage.
// @Tags Scene linkages
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Scene linkage ID"
// @Success 200 {object} v1.ApiResponse[v1.GetSceneLinkageResponse] "Successful response"
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
	v1.HandleSuccess(c, v1.GetSceneLinkageResponse{SceneLinkage: sceneLinkageJSON(*scene)})
}

// CreateSceneLinkage godoc
// @Summary Create scene linkage
// @Schemes
// @Description Creates scene linkage.
// @Tags Scene linkages
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.SceneLinkageRequest true "params"
// @Success 200 {object} v1.ApiResponse[v1.CreateSceneLinkageResponse] "Successful response"
// @Router /scene-linkages [post]
func (h *SceneLinkageHandler) CreateSceneLinkage(c *gin.Context) {
	var req v1.SceneLinkageRequest
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
	v1.HandleSuccess(c, v1.CreateSceneLinkageResponse{SceneLinkage: sceneLinkageJSON(*scene)})
}

// UpdateSceneLinkage godoc
// @Summary Update scene linkage
// @Schemes
// @Description Updates scene linkage.
// @Tags Scene linkages
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Scene linkage ID"
// @Param request body v1.SceneLinkageRequest true "params"
// @Success 200 {object} v1.ApiResponse[v1.UpdateSceneLinkageResponse] "Successful response"
// @Router /scene-linkages/{id} [put]
func (h *SceneLinkageHandler) UpdateSceneLinkage(c *gin.Context) {
	sceneID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	var req v1.SceneLinkageRequest
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
	v1.HandleSuccess(c, v1.UpdateSceneLinkageResponse{SceneLinkage: sceneLinkageJSON(*scene)})
}

// DeleteSceneLinkage godoc
// @Summary Delete scene linkage
// @Schemes
// @Description Deletes scene linkage.
// @Tags Scene linkages
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Scene linkage ID"
// @Success 200 {object} v1.ApiResponse[v1.SceneLinkageSuccessResponse] "Successful response"
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
	v1.HandleSuccess(c, v1.SceneLinkageSuccessResponse{Success: true})
}
