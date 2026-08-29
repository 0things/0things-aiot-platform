package handler

import (
	devicev1 "aiot-backend/api/device/v1"
	devicegroupv1 "aiot-backend/api/device_group/v1"
	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceGroupHandler struct {
	*Handler
	svc service.DeviceGroupServiceInterface
}

func NewDeviceGroupHandler(h *Handler, svc service.DeviceGroupServiceInterface) *DeviceGroupHandler {
	return &DeviceGroupHandler{Handler: h, svc: svc}
}

func deviceGroupJSON(group model.DeviceGroup) devicegroupv1.DeviceGroup {
	return devicegroupv1.DeviceGroup{
		GroupUUID:   group.GroupUUID,
		Name:        group.Name,
		Type:        group.Type,
		Description: group.Description,
		Rule:        group.Rule,
		DeletedAt:   deletedAt(group.DeletedAt),
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
}

// Create godoc
// @Summary 创建设备分组
// @Tags 设备分组
// @Accept json
// @Produce json
// @Param request body devicegroupv1.CreateDeviceGroupRequest true "分组参数"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.DeviceGroup]
// @Router /device-groups [post]
func (h *DeviceGroupHandler) Create(c *gin.Context) {
	var req devicegroupv1.CreateDeviceGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	group, err := h.svc.Create(c, &req)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, deviceGroupJSON(*group))
}

// List godoc
// @Summary 查询设备分组
// @Tags 设备分组
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param search query string false "分组名称"
// @Param type query string false "分组类型"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.ListDeviceGroupsResponse]
// @Router /device-groups [get]
func (h *DeviceGroupHandler) List(c *gin.Context) {
	p, s := page(c, 20)
	groups, total, err := h.svc.List(c, p, s, c.Query("search"), c.Query("type"))
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	items := make([]devicegroupv1.DeviceGroup, len(groups))
	for i := range groups {
		items[i] = deviceGroupJSON(groups[i])
	}
	v1.HandleSuccess(c, devicegroupv1.ListDeviceGroupsResponse{Items: items, Total: total, Page: p, PageSize: s})
}

// Get godoc
// @Summary 查询设备分组详情
// @Tags 设备分组
// @Produce json
// @Param groupUuid path string true "分组 UUID"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.DeviceGroup]
// @Router /device-groups/{groupUuid} [get]
func (h *DeviceGroupHandler) Get(c *gin.Context) {
	group, err := h.svc.Get(c, c.Param("groupUuid"))
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, deviceGroupJSON(*group))
}

// Update godoc
// @Summary 更新设备分组
// @Tags 设备分组
// @Accept json
// @Param groupUuid path string true "分组 UUID"
// @Param request body devicegroupv1.UpdateDeviceGroupRequest true "分组参数"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.DeviceGroup]
// @Router /device-groups/{groupUuid} [put]
func (h *DeviceGroupHandler) Update(c *gin.Context) {
	var req devicegroupv1.UpdateDeviceGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	group, err := h.svc.Update(c, c.Param("groupUuid"), &req)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, deviceGroupJSON(*group))
}

// Delete godoc
// @Summary 删除设备分组
// @Tags 设备分组
// @Param groupUuid path string true "分组 UUID"
// @Success 200 {object} v1.ApiResponse[map[string]bool]
// @Router /device-groups/{groupUuid} [delete]
func (h *DeviceGroupHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c, c.Param("groupUuid")); err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, gin.H{"success": true})
}

// AddDevices godoc
// @Summary 添加分组设备
// @Tags 设备分组
// @Accept json
// @Param groupUuid path string true "分组 UUID"
// @Param request body devicegroupv1.DeviceKeysRequest true "设备 Key 列表"
// @Success 200 {object} v1.ApiResponse[map[string]bool]
// @Router /device-groups/{groupUuid}/devices [post]
func (h *DeviceGroupHandler) AddDevices(c *gin.Context) {
	var req devicegroupv1.DeviceKeysRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	if err := h.svc.AddDevices(c, c.Param("groupUuid"), req.DeviceKeys); err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, gin.H{"success": true})
}

// RemoveDevices godoc
// @Summary 移除分组设备
// @Tags 设备分组
// @Accept json
// @Param groupUuid path string true "分组 UUID"
// @Param request body devicegroupv1.DeviceKeysRequest true "设备 Key 列表"
// @Success 200 {object} v1.ApiResponse[map[string]bool]
// @Router /device-groups/{groupUuid}/devices [delete]
func (h *DeviceGroupHandler) RemoveDevices(c *gin.Context) {
	var req devicegroupv1.DeviceKeysRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	if err := h.svc.RemoveDevices(c, c.Param("groupUuid"), req.DeviceKeys); err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, gin.H{"success": true})
}

// ListDevices godoc
// @Summary 查询分组设备
// @Tags 设备分组
// @Produce json
// @Param groupUuid path string true "分组 UUID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param productKey query string false "产品 Key"
// @Param search query string false "设备 Key 或名称"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.DeviceGroupDevicesResponse]
// @Router /device-groups/{groupUuid}/devices [get]
func (h *DeviceGroupHandler) ListDevices(c *gin.Context) {
	p, s := page(c, 20)
	devices, total, err := h.svc.Devices(c, c.Param("groupUuid"), p, s, c.Query("productKey"), c.Query("search"))
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	items := make([]devicev1.Device, len(devices))
	for i, device := range devices {
		items[i] = deviceJSON(device)
	}
	v1.HandleSuccess(c, devicegroupv1.DeviceGroupDevicesResponse{Devices: items, Total: total, Page: p, PageSize: s})
}

// Preview godoc
// @Summary 预览未保存动态规则
// @Tags 设备分组
// @Accept json
// @Param request body devicegroupv1.PreviewRequest true "动态规则"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.PreviewResponse]
// @Router /device-groups/preview [post]
func (h *DeviceGroupHandler) Preview(c *gin.Context) { h.preview(c) }

// PreviewSaved godoc
// @Summary 预览已保存动态规则
// @Tags 设备分组
// @Accept json
// @Param groupUuid path string true "分组 UUID"
// @Param request body devicegroupv1.PreviewRequest true "动态规则"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.PreviewResponse]
// @Router /device-groups/{groupUuid}/preview [post]
func (h *DeviceGroupHandler) PreviewSaved(c *gin.Context) {
	group, err := h.svc.Get(c, c.Param("groupUuid"))
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	if group.Type != "dynamic" {
		v1.HandleError(c, http.StatusBadRequest, errors.New("only dynamic groups support preview"), nil)
		return
	}
	h.previewRule(c, group.Rule)
}

func (h *DeviceGroupHandler) preview(c *gin.Context) {
	var req devicegroupv1.PreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, http.StatusBadRequest, err, nil)
		return
	}
	h.previewRule(c, req.Rule)
}

func (h *DeviceGroupHandler) previewRule(c *gin.Context, rule string) {
	devices, total, err := h.svc.Preview(c, rule)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	items := make([]devicev1.Device, len(devices))
	for i, device := range devices {
		items[i] = devicev1.Device{
			DeviceKey: device.DeviceKey,
			Name:      device.Name,
			ProductID: device.ProductID,
			Enabled:   device.Enabled,
			CreatedAt: device.CreatedAt,
			UpdatedAt: device.UpdatedAt,
		}
	}
	v1.HandleSuccess(c, devicegroupv1.PreviewResponse{Total: total, Devices: items})
}
