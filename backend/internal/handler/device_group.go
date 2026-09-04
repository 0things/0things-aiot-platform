package handler

import (
	devicegroupv1 "aiot-backend/api/v1"
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
// @Summary Create device group
// @Tags Device groups
// @Accept json
// @Produce json
// @Param request body devicegroupv1.CreateDeviceGroupRequest true "Request payload"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.DeviceGroup] "Successful response"
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
// @Summary List device groups
// @Tags Device groups
// @Produce json
// @Param request query devicegroupv1.ListDeviceGroupsRequest false "Query parameters"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.ListDeviceGroupsResponse] "Successful response"
// @Router /device-groups [get]
func (h *DeviceGroupHandler) List(c *gin.Context) {
	var req devicegroupv1.ListDeviceGroupsRequest
	_ = c.ShouldBindQuery(&req)
	p, s := pageRequest(req.PageRequest, 20)
	groups, total, err := h.svc.List(c, p, s, req.Search, req.Type)
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
// @Summary Get device group
// @Tags Device groups
// @Produce json
// @Param groupUuid path string true "Device group UUID"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.DeviceGroup] "Successful response"
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
// @Summary Update device group
// @Tags Device groups
// @Accept json
// @Param groupUuid path string true "Device group UUID"
// @Param request body devicegroupv1.UpdateDeviceGroupRequest true "Request payload"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.DeviceGroup] "Successful response"
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
// @Summary Delete device group
// @Tags Device groups
// @Param groupUuid path string true "Device group UUID"
// @Success 200 {object} v1.ApiResponse[map[string]bool] "Successful response"
// @Router /device-groups/{groupUuid} [delete]
func (h *DeviceGroupHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c, c.Param("groupUuid")); err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(c, gin.H{"success": true})
}

// AddDevices godoc
// @Summary Add devices to group
// @Tags Device groups
// @Accept json
// @Param groupUuid path string true "Device group UUID"
// @Param request body devicegroupv1.DeviceKeysRequest true "Request payload"
// @Success 200 {object} v1.ApiResponse[map[string]bool] "Successful response"
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
// @Summary Remove devices from group
// @Tags Device groups
// @Accept json
// @Param groupUuid path string true "Device group UUID"
// @Param request body devicegroupv1.DeviceKeysRequest true "Request payload"
// @Success 200 {object} v1.ApiResponse[map[string]bool] "Successful response"
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
// @Summary List devices in group
// @Tags Device groups
// @Produce json
// @Param groupUuid path string true "Device group UUID"
// @Param request query devicegroupv1.ListDeviceGroupDevicesRequest false "Query parameters"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.DeviceGroupDevicesResponse] "Successful response"
// @Router /device-groups/{groupUuid}/devices [get]
func (h *DeviceGroupHandler) ListDevices(c *gin.Context) {
	var req devicegroupv1.ListDeviceGroupDevicesRequest
	_ = c.ShouldBindQuery(&req)
	p, s := pageRequest(req.PageRequest, 20)
	devices, total, err := h.svc.Devices(c, c.Param("groupUuid"), p, s, req.ProductKey, req.Search)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}
	items := make([]devicegroupv1.Device, len(devices))
	for i, device := range devices {
		items[i] = deviceJSON(device)
	}
	v1.HandleSuccess(c, devicegroupv1.DeviceGroupDevicesResponse{Devices: items, Total: total, Page: p, PageSize: s})
}

// Preview godoc
// @Summary Preview device group rule
// @Tags Device groups
// @Accept json
// @Param request body devicegroupv1.PreviewRequest true "Request payload"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.PreviewResponse] "Successful response"
// @Router /device-groups/preview [post]
func (h *DeviceGroupHandler) Preview(c *gin.Context) { h.preview(c) }

// PreviewSaved godoc
// @Summary Preview saved device group rule
// @Tags Device groups
// @Accept json
// @Param groupUuid path string true "Device group UUID"
// @Param request body devicegroupv1.PreviewRequest true "Request payload"
// @Success 200 {object} v1.ApiResponse[devicegroupv1.PreviewResponse] "Successful response"
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
	items := make([]devicegroupv1.Device, len(devices))
	for i, device := range devices {
		items[i] = devicegroupv1.Device{
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
