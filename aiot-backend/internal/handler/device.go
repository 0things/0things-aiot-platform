package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	deviceV1 "0things-backend/api/device/v1"
	"0things-backend/internal/model"
	"0things-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// DeviceHandler deliberately writes the legacy raw response bodies. Device
// clients predate the backend's common response envelope.
type DeviceHandler struct {
	*Handler
	svc    *service.DeviceService
	config *viper.Viper
}

func NewDeviceHandler(h *Handler, svc *service.DeviceService, config *viper.Viper) *DeviceHandler {
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
	_ = json.Unmarshal(history.Desired, &desired)
	_ = json.Unmarshal(history.Reported, &reported)
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

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var req deviceV1.CreateDeviceRequest
	if e := c.ShouldBindJSON(&req); e != nil {
		deviceError(c, e)
		return
	}
	d, e := h.svc.CreateDevice(c, &model.Device{Name: req.Name, ProductID: req.ProductID, Enabled: req.Enabled, Metadata: req.Metadata})
	if e != nil {
		deviceError(c, e)
		return
	}
	c.JSON(200, deviceV1.CreateDeviceResponse{Device: deviceJSON(*d)})
}
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
	c.JSON(200, deviceV1.GetDeviceResponse{Device: deviceJSON(*d)})
}
func (h *DeviceHandler) GetDeviceByKey(c *gin.Context) {
	d, e := h.svc.DeviceByKey(c, c.Param("deviceKey"))
	if e != nil {
		deviceError(c, e)
		return
	}
	c.JSON(200, deviceV1.GetDeviceResponse{Device: deviceJSON(*d)})
}
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
	c.JSON(200, deviceV1.ListDevicesResponse{Devices: out, Total: n, Page: p, PageSize: s})
}
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
	d, e := h.svc.UpdateDevice(c, i, req.Name, req.State, req.Metadata)
	if e != nil {
		deviceError(c, e)
		return
	}
	c.JSON(200, deviceV1.UpdateDeviceResponse{Device: deviceJSON(*d)})
}
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	i, e := id(c)
	if e == nil {
		e = h.svc.DeleteDevice(c, i)
	}
	if e != nil {
		deviceError(c, e)
		return
	}
	c.JSON(200, deviceV1.SuccessResponse{Success: true})
}
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
	c.JSON(200, deviceV1.ActivateDeviceResponse{Device: deviceJSON(*d)})
}
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
	c.JSON(200, deviceV1.SetDeviceEnabledResponse{Device: deviceJSON(*d)})
}
func (h *DeviceHandler) Stats(c *gin.Context) {
	x, e := h.svc.Stats(c)
	if e != nil {
		deviceError(c, e)
		return
	}
	c.JSON(200, deviceV1.DeviceStatisticsResponse{
		TotalDevices: x.TotalDevices, ActivatedDevices: x.ActivatedDevices, OnlineDevices: x.OnlineDevices,
		OfflineDevices: x.OfflineDevices, InactiveDevices: x.InactiveDevices,
	})
}
func (h *DeviceHandler) Telemetry(c *gin.Context) {
	x, e := h.svc.Telemetry(c, c.Param("id"))
	if e != nil {
		deviceError(c, e)
		return
	}
	c.JSON(200, deviceV1.TelemetryResponse{Telemetry: x})
}
func (h *DeviceHandler) MQTT(c *gin.Context) {
	x, e := h.svc.MQTT(c, c.Param("id"))
	if e != nil {
		deviceError(c, e)
		return
	}
	c.JSON(200, deviceV1.MQTTParametersResponse{
		ClientID: x.ClientID, Username: x.Username, MQTTHostURL: x.MQTTHostURL, Password: x.Password, Port: x.Port,
	})
}
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
	c.JSON(200, deviceV1.RestoreDeviceResponse{Device: deviceJSON(*d)})
}

func (h *DeviceHandler) GetTags(c *gin.Context) {
	x, e := h.svc.Tags(c, c.Param("id"))
	if e != nil {
		deviceError(c, e)
		return
	}
	c.JSON(200, deviceV1.ListDeviceTagsResponse{Tags: deviceTagsJSON(x)})
}
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
	c.JSON(200, deviceV1.ListDeviceTagsResponse{Tags: deviceTagsJSON(x)})
}
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
	c.JSON(200, deviceV1.ListDeviceTagsResponse{Tags: deviceTagsJSON(x)})
}
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
	c.JSON(200, deviceV1.SuccessResponse{Success: true})
}
func shadowJSON(x *model.DeviceShadow) deviceV1.Shadow {
	var d, r, m any
	_ = json.Unmarshal(x.Desired, &d)
	_ = json.Unmarshal(x.Reported, &r)
	_ = json.Unmarshal(x.Metadata, &m)
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
func (h *DeviceHandler) GetShadow(c *gin.Context) {
	x, e := h.svc.Shadow(c, c.Param("id"))
	if e != nil {
		deviceError(c, e)
		return
	}
	c.JSON(200, shadowJSON(x))
}
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
	c.JSON(200, shadowJSON(x))
}
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
	c.JSON(200, shadowJSON(x))
}
func (h *DeviceHandler) ClearDesired(c *gin.Context) {
	var req deviceV1.ClearDesiredShadowRequest
	_ = c.ShouldBindJSON(&req)
	x, e := h.svc.MutateShadow(c, c.Param("id"), req.Version, "app", nil, nil, true)
	if e != nil {
		deviceError(c, e)
		return
	}
	c.JSON(200, shadowJSON(x))
}
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
	c.JSON(200, deviceV1.ListDeviceShadowHistoryResponse{History: items})
}

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
	c.JSON(200, deviceV1.SimulatePushResponse{PushRecordID: strconv.FormatInt(record.ID, 10), Timestamp: record.CreatedAt.UnixMilli(), Status: record.Status, Message: "success"})
}
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
	c.JSON(200, deviceV1.ListPushRecordsResponse{Records: items, Total: total, Page: p, PageSize: s})
}
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
	c.JSON(200, deviceV1.GetPushRecordResponse{Record: pushRecordJSON(*record)})
}
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
	c.JSON(200, deviceV1.ClearPushRecordsResponse{DeletedCount: deletedCount, Success: true})
}

func (h *DeviceHandler) BatchTemplate(c *gin.Context) {
	b, e := h.svc.BatchTemplate()
	if e != nil {
		deviceError(c, e)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=device_import_template.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", b)
}
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
	c.JSON(200, deviceV1.BatchUploadDevicesResponse{SuccessCount: n, FailureCount: len(items), Errors: items})
}
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
	c.JSON(200, deviceV1.MockKafkaResponse{Success: true, Message: "Message sent successfully to topic: " + req.Topic})
}
