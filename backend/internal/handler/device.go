package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	deviceV1 "aiot-backend/api/device/v1"
	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// DeviceHandler deliberately writes the legacy raw response bodies. Device
// clients predate the backend's common response envelope.
type DeviceHandler struct {
	*Handler
	svc    service.DeviceServiceInterface
	config *viper.Viper
}

func NewDeviceHandler(h *Handler, svc service.DeviceServiceInterface, config *viper.Viper) *DeviceHandler {
	return &DeviceHandler{Handler: h, svc: svc, config: config}
}

func deviceJSON(device model.Device) deviceV1.Device {
	return deviceV1.Device{
		ID:              device.ID,
		DeviceKey:       device.DeviceKey,
		Name:            device.Name,
		ProductID:       device.ProductID,
		ProductKey:      device.Product.ProductKey,
		ProductName:     device.Product.Name,
		State:           device.State.State,
		Enabled:         device.Enabled,
		LastOnlineTime:  device.State.LastOnlineTime,
		LastOfflineTime: device.State.LastOfflineTime,
		Metadata:        fmt.Sprint(raw(device.Metadata)),
		CreatedAt:       device.CreatedAt,
		UpdatedAt:       device.UpdatedAt,
		DeletedAt:       deletedAt(device.DeletedAt),
	}
}

func deviceTagJSON(tag model.DeviceTag) deviceV1.DeviceTag {
	return deviceV1.DeviceTag{
		ID: tag.ID, DeviceID: tag.DeviceID, Key: tag.Key, Value: tag.Value, Source: tag.Source,
		CreatedAt: tag.CreatedAt, UpdatedAt: tag.UpdatedAt, DeletedAt: deletedAt(tag.DeletedAt),
	}
}

func deviceTagsJSON(tags []model.DeviceTag) []deviceV1.DeviceTag {
	items := make([]deviceV1.DeviceTag, len(tags))
	for i, tag := range tags {
		items[i] = deviceTagJSON(tag)
	}
	return items
}

func shadowHistoryJSON(history model.DeviceShadowHistory) deviceV1.DeviceShadowHistory {
	var desired, reported any
	_ = json.Unmarshal([]byte(history.Desired), &desired)
	_ = json.Unmarshal([]byte(history.Reported), &reported)
	return deviceV1.DeviceShadowHistory{
		ID: history.ID, DeviceID: history.DeviceID, Version: history.Version, Source: history.Source,
		Desired: desired, Reported: reported, CreatedAt: history.CreatedAt,
	}
}

func pushRecordJSON(record model.DevicePushRecord) deviceV1.PushRecord {
	return deviceV1.PushRecord{
		ID: record.ID, DeviceID: record.DeviceID, OperationType: record.OperationType, OperationName: record.OperationName,
		Payload: record.Payload, Status: record.Status, ErrorMessage: record.ErrorMessage, CreatedBy: record.CreatedBy,
		CreatedAt: record.CreatedAt, UpdatedAt: record.UpdatedAt,
	}
}

// CreateDevice godoc
// @Summary 创建设备
// @Schemes
// @Description 创建一个新的设备
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body deviceV1.CreateDeviceRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.CreateDeviceResponse]
// @Router /devices [post]
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var req deviceV1.CreateDeviceRequest
	if e := c.ShouldBindJSON(&req); e != nil {
		deviceError(c, e)
		return
	}
	d, e := h.svc.CreateDevice(c, &model.Device{Name: req.Name, ProductID: req.ProductID, Enabled: req.Enabled, Metadata: string(req.Metadata)})
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.CreateDeviceResponse{Device: deviceJSON(*d)})
}
// GetDevice godoc
// @Summary 通过 ID 获取设备
// @Schemes
// @Description 通过设备 ID 获取设备详情
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Success 200 {object} v1.ApiResponse[deviceV1.GetDeviceResponse]
// @Router /devices/{id} [get]
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	i, e := id(c)
	if e != nil {
		deviceError(c, e)
		return
	}
	d, e := h.svc.Device(c, i)
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.GetDeviceResponse{Device: deviceJSON(*d)})
}
// GetDeviceByKey godoc
// @Summary 通过 deviceKey 获取设备
// @Schemes
// @Description 通过设备 Key 获取设备详情
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "设备 Key"
// @Success 200 {object} v1.ApiResponse[deviceV1.GetDeviceResponse]
// @Router /devices/key/{deviceKey} [get]
func (h *DeviceHandler) GetDeviceByKey(c *gin.Context) {
	d, e := h.svc.DeviceByKey(c, c.Param("deviceKey"))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.GetDeviceResponse{Device: deviceJSON(*d)})
}
// ListDevices godoc
// @Summary 获取设备列表
// @Schemes
// @Description 分页获取设备列表，支持按 productId、states、enabled、searchText 过滤
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param productId query int false "产品 ID"
// @Param states query []string false "设备状态"
// @Param enabled query bool false "是否启用"
// @Param searchText query string false "搜索关键字"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListDevicesResponse]
// @Router /devices [get]
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	p, s := page(c, 10)
	pid, _ := strconv.ParseInt(c.Query("productId"), 10, 64)
	var enabled *bool
	if x, ok := c.GetQuery("enabled"); ok {
		v, _ := strconv.ParseBool(x)
		enabled = &v
	}
	d, n, e := h.svc.ListDevices(c, p, s, pid, c.QueryArray("states"), enabled, c.Query("searchText"))
	if e != nil {
		deviceError(c, e)
		return
	}
	out := make([]deviceV1.Device, 0, len(d))
	for _, v := range d {
		out = append(out, deviceJSON(v))
	}
	v1.HandleSuccess(c, deviceV1.ListDevicesResponse{Devices: out, Total: n, Page: p, PageSize: s})
}
// UpdateDevice godoc
// @Summary 更新设备
// @Schemes
// @Description 通过设备 ID 更新设备
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Param request body deviceV1.UpdateDeviceRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.UpdateDeviceResponse]
// @Router /devices/{id} [put]
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	i, e := id(c)
	if e != nil {
		deviceError(c, e)
		return
	}
	var req deviceV1.UpdateDeviceRequest
	if e = c.ShouldBindJSON(&req); e != nil {
		deviceError(c, e)
		return
	}
	d, e := h.svc.UpdateDevice(c, i, req.Name, req.State, string(req.Metadata))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.UpdateDeviceResponse{Device: deviceJSON(*d)})
}
// DeleteDevice godoc
// @Summary 删除设备
// @Schemes
// @Description 通过设备 ID 软删除设备
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Success 200 {object} v1.ApiResponse[deviceV1.SuccessResponse]
// @Router /devices/{id} [delete]
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	i, e := id(c)
	if e == nil {
		e = h.svc.DeleteDevice(c, i)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.SuccessResponse{Success: true})
}
// Activate godoc
// @Summary 激活设备
// @Schemes
// @Description 通过设备 ID 激活设备
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Success 200 {object} v1.ApiResponse[deviceV1.ActivateDeviceResponse]
// @Router /devices/{id}/activate [post]
func (h *DeviceHandler) Activate(c *gin.Context) {
	i, e := id(c)
	if e != nil {
		deviceError(c, e)
		return
	}
	d, e := h.svc.Activate(c, i)
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.ActivateDeviceResponse{Device: deviceJSON(*d)})
}
// Enabled godoc
// @Summary 设置设备启用/禁用
// @Schemes
// @Description 通过设备 ID 启用或禁用设备
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Param request body deviceV1.SetDeviceEnabledRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.SetDeviceEnabledResponse]
// @Router /devices/{id}/enabled [post]
func (h *DeviceHandler) Enabled(c *gin.Context) {
	i, e := id(c)
	var req deviceV1.SetDeviceEnabledRequest
	if e == nil {
		e = c.ShouldBindJSON(&req)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	d, e := h.svc.SetEnabled(c, i, req.Enabled)
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.SetDeviceEnabledResponse{Device: deviceJSON(*d)})
}
// Stats godoc
// @Summary 设备统计
// @Schemes
// @Description 获取设备统计信息（总数、已激活、在线、离线、未激活）
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.ApiResponse[deviceV1.DeviceStatisticsResponse]
// @Router /device-statistics [get]
func (h *DeviceHandler) Stats(c *gin.Context) {
	x, e := h.svc.Stats(c)
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.DeviceStatisticsResponse{
		TotalDevices: x.TotalDevices, ActivatedDevices: x.ActivatedDevices, OnlineDevices: x.OnlineDevices,
		OfflineDevices: x.OfflineDevices, InactiveDevices: x.InactiveDevices,
	})
}
// Telemetry godoc
// @Summary 获取设备遥测数据
// @Schemes
// @Description 获取设备的遥测（Telemetry）数据
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Success 200 {object} v1.ApiResponse[deviceV1.TelemetryResponse]
// @Router /devices/{id}/telemetry [get]
func (h *DeviceHandler) Telemetry(c *gin.Context) {
	x, e := h.svc.Telemetry(c, c.Param("id"))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.TelemetryResponse{Telemetry: x})
}
// MQTT godoc
// @Summary 获取 MQTT 连接参数
// @Schemes
// @Description 获取设备 MQTT 连接的 ClientID、Username、Password、HostURL、Port
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param device_key path string true "设备 DeviceKey"
// @Success 200 {object} v1.ApiResponse[deviceV1.MQTTParametersResponse]
// @Router /devices/{device_key}/mqtt-parameters [get]
func (h *DeviceHandler) MQTT(c *gin.Context) {
	x, e := h.svc.MQTT(c, c.Param("id"))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.MQTTParametersResponse{
		ClientID: x.ClientID, Username: x.Username, MQTTHostURL: x.MQTTHostURL, Password: x.Password, Port: x.Port,
	})
}
// Restore godoc
// @Summary 恢复已删除的设备
// @Schemes
// @Description 通过设备 ID 恢复软删除的设备
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Success 200 {object} v1.ApiResponse[deviceV1.RestoreDeviceResponse]
// @Router /devices/{id}/restore [post]
func (h *DeviceHandler) Restore(c *gin.Context) {
	i, e := id(c)
	if e != nil {
		deviceError(c, e)
		return
	}
	d, e := h.svc.RestoreDevice(c, i)
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.RestoreDeviceResponse{Device: deviceJSON(*d)})
}

// GetTags godoc
// @Summary 获取设备标签
// @Schemes
// @Description 获取设备的所有标签
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListDeviceTagsResponse]
// @Router /devices/{id}/tags [get]
func (h *DeviceHandler) GetTags(c *gin.Context) {
	x, e := h.svc.Tags(c, c.Param("id"))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.ListDeviceTagsResponse{Tags: deviceTagsJSON(x)})
}
// PutTags godoc
// @Summary 全量覆盖设置设备标签
// @Schemes
// @Description 覆盖式设置设备的标签集合（PUT 语义）
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Param request body deviceV1.SetDeviceTagsRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListDeviceTagsResponse]
// @Router /devices/{id}/tags [put]
func (h *DeviceHandler) PutTags(c *gin.Context) {
	var req deviceV1.SetDeviceTagsRequest
	e := c.ShouldBindJSON(&req)
	var x []model.DeviceTag
	if e == nil {
		x, e = h.svc.SetTags(c, c.Param("id"), req.Tags, true)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.ListDeviceTagsResponse{Tags: deviceTagsJSON(x)})
}
// PostTags godoc
// @Summary 增量添加设备标签
// @Schemes
// @Description 增量式添加设备的标签（POST 语义，保留已存在的标签）
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Param request body deviceV1.SetDeviceTagsRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListDeviceTagsResponse]
// @Router /devices/{id}/tags [post]
func (h *DeviceHandler) PostTags(c *gin.Context) {
	var req deviceV1.SetDeviceTagsRequest
	e := c.ShouldBindJSON(&req)
	var x []model.DeviceTag
	if e == nil {
		x, e = h.svc.SetTags(c, c.Param("id"), req.Tags, false)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.ListDeviceTagsResponse{Tags: deviceTagsJSON(x)})
}
// DeleteTags godoc
// @Summary 删除设备标签
// @Schemes
// @Description 删除设备指定的标签键
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Param request body deviceV1.DeleteDeviceTagsRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.SuccessResponse]
// @Router /devices/{id}/tags [delete]
func (h *DeviceHandler) DeleteTags(c *gin.Context) {
	var req deviceV1.DeleteDeviceTagsRequest
	e := c.ShouldBindJSON(&req)
	if e == nil {
		e = h.svc.RemoveTags(c, c.Param("id"), req.Keys)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.SuccessResponse{Success: true})
}
func shadowJSON(x *model.DeviceShadow) deviceV1.Shadow {
	var d, r, m any
	_ = json.Unmarshal([]byte(x.Desired), &d)
	_ = json.Unmarshal([]byte(x.Reported), &r)
	_ = json.Unmarshal([]byte(x.Metadata), &m)
	dm := map[string]any{}
	if ds, ok := d.(map[string]any); ok {
		rs, _ := r.(map[string]any)
		for k, v := range ds {
			if fmt.Sprint(rs[k]) != fmt.Sprint(v) {
				dm[k] = v
			}
		}
	}
	return deviceV1.Shadow{Desired: d, Reported: r, Delta: dm, Metadata: m, Version: x.Version, UpdatedAt: x.UpdatedAt}
}
// GetShadow godoc
// @Summary 获取设备影子
// @Schemes
// @Description 获取设备影子（Desired、Reported、Delta、Metadata、Version）
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Success 200 {object} v1.ApiResponse[deviceV1.Shadow]
// @Router /devices/{id}/shadow [get]
func (h *DeviceHandler) GetShadow(c *gin.Context) {
	x, e := h.svc.Shadow(c, c.Param("id"))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, shadowJSON(x))
}
// Desired godoc
// @Summary 更新设备影子期望值
// @Schemes
// @Description 由应用侧更新设备的 Desired 影子
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Param request body deviceV1.UpdateDesiredShadowRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.Shadow]
// @Router /devices/{id}/shadow/desired [put]
func (h *DeviceHandler) Desired(c *gin.Context) {
	var req deviceV1.UpdateDesiredShadowRequest
	e := c.ShouldBindJSON(&req)
	var x *model.DeviceShadow
	if e == nil {
		x, e = h.svc.MutateShadow(c, c.Param("id"), req.Version, "app", &req.Desired, nil, false)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, shadowJSON(x))
}
// Reported godoc
// @Summary 更新设备影子上报值
// @Schemes
// @Description 由设备侧更新 Reported 影子
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Param request body deviceV1.UpdateReportedShadowRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.Shadow]
// @Router /devices/{id}/shadow/reported [put]
func (h *DeviceHandler) Reported(c *gin.Context) {
	var req deviceV1.UpdateReportedShadowRequest
	e := c.ShouldBindJSON(&req)
	var x *model.DeviceShadow
	if e == nil {
		x, e = h.svc.MutateShadow(c, c.Param("id"), req.Version, "device", nil, &req.Reported, false)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, shadowJSON(x))
}
// ClearDesired godoc
// @Summary 清空设备影子期望值
// @Schemes
// @Description 清空设备影子的 Desired 部分
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Param request body deviceV1.ClearDesiredShadowRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.Shadow]
// @Router /devices/{id}/shadow/desired [delete]
func (h *DeviceHandler) ClearDesired(c *gin.Context) {
	var req deviceV1.ClearDesiredShadowRequest
	_ = c.ShouldBindJSON(&req)
	x, e := h.svc.MutateShadow(c, c.Param("id"), req.Version, "app", nil, nil, true)
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, shadowJSON(x))
}
// History godoc
// @Summary 获取设备影子历史
// @Schemes
// @Description 获取设备影子的变更历史记录
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListDeviceShadowHistoryResponse]
// @Router /devices/{id}/shadow/history [get]
func (h *DeviceHandler) History(c *gin.Context) {
	x, e := h.svc.ShadowHistory(c, c.Param("id"))
	if e != nil {
		deviceError(c, e)
		return
	}
	items := make([]deviceV1.DeviceShadowHistory, len(x))
	for i, history := range x {
		items[i] = shadowHistoryJSON(history)
	}
	v1.HandleSuccess(c, deviceV1.ListDeviceShadowHistoryResponse{History: items})
}

// SimulatePush godoc
// @Summary 模拟下行推送
// @Schemes
// @Description 模拟服务端向设备发起一次下行推送
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Param request body deviceV1.SimulatePushRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.SimulatePushResponse]
// @Router /devices/{id}/simulate-push [post]
func (h *DeviceHandler) SimulatePush(c *gin.Context) {
	var req deviceV1.SimulatePushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	record, err := h.svc.SimulatePush(c, c.Param("id"), req.Payload, c.GetHeader("X-User-ID"))
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, deviceV1.SimulatePushResponse{PushRecordID: strconv.FormatInt(record.ID, 10), Timestamp: record.CreatedAt.UnixMilli(), Status: record.Status, Message: "success"})
}
// PushRecords godoc
// @Summary 获取设备下行推送记录列表
// @Schemes
// @Description 分页获取设备的下行推送记录
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param operationType query string false "操作类型"
// @Param status query string false "状态"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListPushRecordsResponse]
// @Router /devices/{id}/push-records [get]
func (h *DeviceHandler) PushRecords(c *gin.Context) {
	p, s := page(c, 20)
	records, total, err := h.svc.ListPushRecords(c, c.Param("id"), p, s, c.Query("operationType"), c.Query("status"))
	if err != nil {
		deviceError(c, err)
		return
	}
	items := make([]deviceV1.PushRecord, len(records))
	for i, record := range records {
		items[i] = pushRecordJSON(record)
	}
	v1.HandleSuccess(c, deviceV1.ListPushRecordsResponse{Records: items, Total: total, Page: p, PageSize: s})
}
// PushRecord godoc
// @Summary 获取设备下行推送记录
// @Schemes
// @Description 通过推送记录 ID 获取单个下行推送记录
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param pushRecordId path int true "推送记录 ID"
// @Success 200 {object} v1.ApiResponse[deviceV1.GetPushRecordResponse]
// @Router /devices/push-records/{pushRecordId} [get]
func (h *DeviceHandler) PushRecord(c *gin.Context) {
	pushRecordID, err := strconv.ParseInt(c.Param("pushRecordId"), 10, 64)
	if err != nil {
		deviceError(c, err)
		return
	}
	record, err := h.svc.PushRecord(c, pushRecordID)
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, deviceV1.GetPushRecordResponse{Record: pushRecordJSON(*record)})
}
// ClearPushRecords godoc
// @Summary 清理设备下行推送记录
// @Schemes
// @Description 清理设备指定时间戳之前的下行推送记录
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设备 ID"
// @Param beforeTimestamp query int false "清理此时间戳之前的记录（毫秒）"
// @Success 200 {object} v1.ApiResponse[deviceV1.ClearPushRecordsResponse]
// @Router /devices/{id}/push-records [delete]
func (h *DeviceHandler) ClearPushRecords(c *gin.Context) {
	var before *time.Time
	if v := c.Query("beforeTimestamp"); v != "" {
		if milliseconds, _ := strconv.ParseInt(v, 10, 64); milliseconds > 0 {
			value := time.UnixMilli(milliseconds)
			before = &value
		}
	}
	deletedCount, err := h.svc.ClearPushRecords(c, c.Param("id"), before)
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, deviceV1.ClearPushRecordsResponse{DeletedCount: deletedCount, Success: true})
}

// BatchTemplate godoc
// @Summary 下载设备批量导入模板
// @Schemes
// @Description 下载设备批量导入的 Excel 模板
// @Tags 设备模块
// @Produce octet-stream
// @Security Bearer
// @Success 200 {file} file
// @Router /devices/batch/template [get]
func (h *DeviceHandler) BatchTemplate(c *gin.Context) {
	b, e := h.svc.BatchTemplate()
	if e != nil {
		deviceError(c, e)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=device_import_template.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", b)
}
// BatchUpload godoc
// @Summary 批量上传设备
// @Schemes
// @Description 通过 Excel 文件批量导入设备
// @Tags 设备模块
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "Excel 模板文件"
// @Success 200 {object} v1.ApiResponse[deviceV1.BatchUploadDevicesResponse]
// @Router /devices/batch/upload [post]
func (h *DeviceHandler) BatchUpload(c *gin.Context) {
	f, _, e := c.Request.FormFile("file")
	if e != nil {
		deviceError(c, e)
		return
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, 10<<20))
	if e != nil {
		deviceError(c, e)
		return
	}
	n, errs, e := h.svc.BatchCreate(c, b)
	if e != nil {
		deviceError(c, e)
		return
	}
	items := make([]deviceV1.BatchUploadError, len(errs))
	for i, item := range errs {
		items[i] = deviceV1.BatchUploadError{Row: item.Row, ProductKey: item.ProductKey, DeviceName: item.DeviceName, Error: item.Error}
	}
	v1.HandleSuccess(c, deviceV1.BatchUploadDevicesResponse{SuccessCount: n, FailureCount: len(items), Errors: items})
}
// MockKafka godoc
// @Summary 模拟向 Kafka 发送消息
// @Schemes
// @Description 向指定 Kafka topic 发送一条测试消息
// @Tags 设备模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body deviceV1.MockKafkaRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.MockKafkaResponse]
// @Router /devices/mock-kafka [post]
func (h *DeviceHandler) MockKafka(c *gin.Context) {
	var req deviceV1.MockKafkaRequest
	e := c.ShouldBindJSON(&req)
	if e == nil {
		e = h.svc.MockKafka(c, h.config.GetStringSlice("data.kafka.device.brokers"), req.Topic, req.Data)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.MockKafkaResponse{Success: true, Message: "Message sent successfully to topic: " + req.Topic})
}
