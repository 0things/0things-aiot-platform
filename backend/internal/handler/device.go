package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	deviceV1 "aiot-backend/api/v1"
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
// @Summary Create device
// @Schemes
// @Description Creates device.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body deviceV1.CreateDeviceRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.CreateDeviceResponse] "Successful response"
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
// @Summary Get device
// @Schemes
// @Description Returns device.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Success 200 {object} v1.ApiResponse[deviceV1.GetDeviceResponse] "Successful response"
// @Router /devices/{deviceKey} [get]
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	d, e := h.svc.DeviceByKey(c, c.Param("deviceKey"))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.GetDeviceResponse{Device: deviceJSON(*d)})
}

// GetDeviceByKey preserves the internal entry point; the public route uses /devices/:deviceKey.
func (h *DeviceHandler) GetDeviceByKey(c *gin.Context) {
	h.GetDevice(c)
}

// ListDevices godoc
// @Summary List devices
// @Schemes
// @Description Lists devices with pagination and optional product, state, enabled, and text filters.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query deviceV1.ListDevicesRequest false "Query parameters"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListDevicesResponse] "Successful response"
// @Router /devices [get]
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	var req deviceV1.ListDevicesRequest
	_ = c.ShouldBindQuery(&req)
	p, s := pageRequest(req.PageRequest, 10)
	d, n, e := h.svc.ListDevices(c, p, s, req.ProductID, req.States, req.Enabled, req.SearchText)
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
// @Summary Update device
// @Schemes
// @Description Updates device.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param request body deviceV1.UpdateDeviceRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.UpdateDeviceResponse] "Successful response"
// @Router /devices/{deviceKey} [put]
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	var req deviceV1.UpdateDeviceRequest
	if e := c.ShouldBindJSON(&req); e != nil {
		deviceError(c, e)
		return
	}
	d, e := h.svc.UpdateDeviceByKey(c, c.Param("deviceKey"), req.Name, req.State, string(req.Metadata))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.UpdateDeviceResponse{Device: deviceJSON(*d)})
}

// DeleteDevice godoc
// @Summary Delete Device
// @Schemes
// @Description Deletes Device.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Success 200 {object} v1.ApiResponse[deviceV1.DeviceSuccessResponse] "Successful response"
// @Router /devices/{deviceKey} [delete]
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	e := h.svc.DeleteDeviceByKey(c, c.Param("deviceKey"))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.DeviceSuccessResponse{Success: true})
}

// Activate godoc
// @Summary Activate device
// @Schemes
// @Description Activates device.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Success 200 {object} v1.ApiResponse[deviceV1.ActivateDeviceResponse] "Successful response"
// @Router /devices/{deviceKey}/activate [post]
func (h *DeviceHandler) Activate(c *gin.Context) {
	d, e := h.svc.ActivateByKey(c, c.Param("deviceKey"))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.ActivateDeviceResponse{Device: deviceJSON(*d)})
}

// Enabled godoc
// @Summary Set device enabled state
// @Schemes
// @Description Sets device enabled state.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param request body deviceV1.SetDeviceEnabledRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.SetDeviceEnabledResponse] "Successful response"
// @Router /devices/{deviceKey}/enabled [post]
func (h *DeviceHandler) Enabled(c *gin.Context) {
	var req deviceV1.SetDeviceEnabledRequest
	e := c.ShouldBindJSON(&req)
	if e != nil {
		deviceError(c, e)
		return
	}
	d, e := h.svc.SetEnabledByKey(c, c.Param("deviceKey"), req.Enabled)
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.SetDeviceEnabledResponse{Device: deviceJSON(*d)})
}

// Stats godoc
// @Summary Get device statistics
// @Schemes
// @Description Returns device statistics.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.ApiResponse[deviceV1.DeviceStatisticsResponse] "Successful response"
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
// @Summary Get device telemetry
// @Schemes
// @Description Returns device telemetry.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Success 200 {object} v1.ApiResponse[deviceV1.TelemetryResponse] "Successful response"
// @Router /devices/{deviceKey}/telemetry [get]
func (h *DeviceHandler) Telemetry(c *gin.Context) {
	x, e := h.svc.Telemetry(c, c.Param("deviceKey"))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.TelemetryResponse{Telemetry: x})
}

// GetTags godoc
// @Summary List device tags
// @Schemes
// @Description Lists device tags.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListDeviceTagsResponse] "Successful response"
// @Router /devices/{deviceKey}/tags [get]
func (h *DeviceHandler) GetTags(c *gin.Context) {
	x, e := h.svc.Tags(c, c.Param("deviceKey"))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.ListDeviceTagsResponse{Tags: deviceTagsJSON(x)})
}

// PutTags godoc
// @Summary Replace device tags
// @Schemes
// @Description Replaces device tags.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param request body deviceV1.SetDeviceTagsRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListDeviceTagsResponse] "Successful response"
// @Router /devices/{deviceKey}/tags [put]
func (h *DeviceHandler) PutTags(c *gin.Context) {
	var req deviceV1.SetDeviceTagsRequest
	e := c.ShouldBindJSON(&req)
	var x []model.DeviceTag
	if e == nil {
		x, e = h.svc.SetTags(c, c.Param("deviceKey"), req.Tags, true)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.ListDeviceTagsResponse{Tags: deviceTagsJSON(x)})
}

// PostTags godoc
// @Summary Add device tags
// @Schemes
// @Description Adds device tags.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param request body deviceV1.SetDeviceTagsRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListDeviceTagsResponse] "Successful response"
// @Router /devices/{deviceKey}/tags [post]
func (h *DeviceHandler) PostTags(c *gin.Context) {
	var req deviceV1.SetDeviceTagsRequest
	e := c.ShouldBindJSON(&req)
	var x []model.DeviceTag
	if e == nil {
		x, e = h.svc.SetTags(c, c.Param("deviceKey"), req.Tags, false)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.ListDeviceTagsResponse{Tags: deviceTagsJSON(x)})
}

// DeleteTags godoc
// @Summary Delete device tags
// @Schemes
// @Description Deletes device tags.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param request body deviceV1.DeleteDeviceTagsRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.DeviceSuccessResponse] "Successful response"
// @Router /devices/{deviceKey}/tags [delete]
func (h *DeviceHandler) DeleteTags(c *gin.Context) {
	var req deviceV1.DeleteDeviceTagsRequest
	e := c.ShouldBindJSON(&req)
	if e == nil {
		e = h.svc.RemoveTags(c, c.Param("deviceKey"), req.Keys)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, deviceV1.DeviceSuccessResponse{Success: true})
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
	// When metadata is absent, synthesize the standard per-property timestamp structure.
	if m == nil || m == "" {
		metaMap := map[string]any{
			"desired":  map[string]any{},
			"reported": map[string]any{},
		}
		nowTs := x.UpdatedAt.Unix()
		if ds, ok := d.(map[string]any); ok {
			dMeta := map[string]any{}
			for k := range ds {
				dMeta[k] = map[string]any{"timestamp": nowTs}
			}
			metaMap["desired"] = dMeta
		}
		if rs, ok := r.(map[string]any); ok {
			rMeta := map[string]any{}
			for k := range rs {
				rMeta[k] = map[string]any{"timestamp": nowTs}
			}
			metaMap["reported"] = rMeta
		}
		m = metaMap
	}

	state := deviceV1.ShadowState{
		Desired:  d,
		Reported: r,
		Delta:    dm,
	}

	return deviceV1.Shadow{
		State:     state,
		Desired:   d,
		Reported:  r,
		Delta:     dm,
		Metadata:  m,
		Version:   x.Version,
		Timestamp: x.UpdatedAt.Unix(),
		UpdatedAt: x.UpdatedAt,
	}
}

// GetShadow godoc
// @Summary Get device shadow
// @Schemes
// @Description Returns device shadow.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Success 200 {object} v1.ApiResponse[deviceV1.Shadow] "Successful response"
// @Router /devices/{deviceKey}/shadow [get]
func (h *DeviceHandler) GetShadow(c *gin.Context) {
	x, e := h.svc.Shadow(c, c.Param("deviceKey"))
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, shadowJSON(x))
}

// Desired godoc
// @Summary Update desired device shadow
// @Schemes
// @Description Updates desired device shadow.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param request body deviceV1.UpdateDesiredShadowRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.Shadow] "Successful response"
// @Router /devices/{deviceKey}/shadow/desired [put]
func (h *DeviceHandler) Desired(c *gin.Context) {
	var req deviceV1.UpdateDesiredShadowRequest
	e := c.ShouldBindJSON(&req)
	var x *model.DeviceShadow
	if e == nil {
		x, e = h.svc.MutateShadow(c, c.Param("deviceKey"), req.Version, "app", &req.Desired, nil, false)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, shadowJSON(x))
}

// ClearDesired godoc
// @Summary Clear desired device shadow
// @Schemes
// @Description Clears desired device shadow.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param request body deviceV1.ClearDesiredShadowRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.Shadow] "Successful response"
// @Router /devices/{deviceKey}/shadow/desired [delete]
func (h *DeviceHandler) ClearDesired(c *gin.Context) {
	var req deviceV1.ClearDesiredShadowRequest
	_ = c.ShouldBindJSON(&req)
	x, e := h.svc.MutateShadow(c, c.Param("deviceKey"), req.Version, "app", nil, nil, true)
	if e != nil {
		deviceError(c, e)
		return
	}
	v1.HandleSuccess(c, shadowJSON(x))
}

// History godoc
// @Summary List device shadow history
// @Schemes
// @Description Lists device shadow history.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListDeviceShadowHistoryResponse] "Successful response"
// @Router /devices/{deviceKey}/shadow/history [get]
func (h *DeviceHandler) History(c *gin.Context) {
	x, e := h.svc.ShadowHistory(c, c.Param("deviceKey"))
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
// @Summary Simulate downstream push
// @Schemes
// @Description Simulates downstream push.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param request body deviceV1.SimulatePushRequest true "params"
// @Success 200 {object} v1.ApiResponse[deviceV1.SimulatePushResponse] "Successful response"
// @Router /devices/{deviceKey}/simulate-push [post]
func (h *DeviceHandler) SimulatePush(c *gin.Context) {
	var req deviceV1.SimulatePushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		deviceError(c, err)
		return
	}
	record, err := h.svc.SimulatePush(c, c.Param("deviceKey"), req.Payload, c.GetHeader("X-User-ID"))
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, deviceV1.SimulatePushResponse{PushRecordID: strconv.FormatInt(record.ID, 10), Timestamp: record.CreatedAt.UnixMilli(), Status: record.Status, Message: "success"})
}

// PushRecords godoc
// @Summary List downstream push records
// @Schemes
// @Description Lists downstream push records.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param request query deviceV1.ListPushRecordsRequest false "Query parameters"
// @Success 200 {object} v1.ApiResponse[deviceV1.ListPushRecordsResponse] "Successful response"
// @Router /devices/{deviceKey}/push-records [get]
func (h *DeviceHandler) PushRecords(c *gin.Context) {
	var req deviceV1.ListPushRecordsRequest
	_ = c.ShouldBindQuery(&req)
	p, s := pageRequest(req.PageRequest, 20)
	records, total, err := h.svc.ListPushRecords(c, c.Param("deviceKey"), p, s, req.OperationType, req.Status)
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
// @Summary Get downstream push record
// @Schemes
// @Description Returns downstream push record.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param pushRecordId path int true "Push record ID"
// @Success 200 {object} v1.ApiResponse[deviceV1.GetPushRecordResponse] "Successful response"
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
// @Summary Clear downstream push records
// @Schemes
// @Description Clears downstream push records.
// @Tags Devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param beforeTimestamp query int false "Cutoff timestamp"
// @Success 200 {object} v1.ApiResponse[deviceV1.ClearPushRecordsResponse] "Successful response"
// @Router /devices/{deviceKey}/push-records [delete]
func (h *DeviceHandler) ClearPushRecords(c *gin.Context) {
	var before *time.Time
	if v := c.Query("beforeTimestamp"); v != "" {
		if milliseconds, _ := strconv.ParseInt(v, 10, 64); milliseconds > 0 {
			value := time.UnixMilli(milliseconds)
			before = &value
		}
	}
	deletedCount, err := h.svc.ClearPushRecords(c, c.Param("deviceKey"), before)
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, deviceV1.ClearPushRecordsResponse{DeletedCount: deletedCount, Success: true})
}

// BatchTemplate godoc
// @Summary Download device import template
// @Schemes
// @Description Downloads device import template.
// @Tags Devices
// @Produce octet-stream
// @Security Bearer
// @Success 200 {file} file "Successful response"
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
// @Summary Import devices in bulk
// @Schemes
// @Description Imports devices in bulk.
// @Tags Devices
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "Upload file"
// @Success 200 {object} v1.ApiResponse[deviceV1.BatchUploadDevicesResponse] "Successful response"
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
