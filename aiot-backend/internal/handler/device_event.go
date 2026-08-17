package handler

import (
	"net/http"
	"time"

	eventV1 "0things-backend/api/event/v1"
	v1 "0things-backend/api/v1"
	"0things-backend/internal/model"
	"0things-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type DeviceEventHandler struct { *Handler; svc service.DeviceEventServiceInterface }

func NewDeviceEventHandler(h *Handler, svc service.DeviceEventServiceInterface) *DeviceEventHandler { return &DeviceEventHandler{Handler: h, svc: svc} }

func deviceEventResponse(event model.DeviceEvent) eventV1.DeviceEvent {
	return eventV1.DeviceEvent{ID: event.ID, DeviceKey: event.DeviceKey, DeviceName: event.DeviceName, EventType: event.EventType, EventAt: event.EventAt, Data: string(event.Data)}
}

// ListDeviceEvents godoc
// @Summary 获取设备事件列表
// @Description 分页获取物模型设备事件，支持关键词、设备、事件类型与时间范围筛选
// @Tags 事件管理
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键词"
// @Param device_key query string false "设备 Key"
// @Param event_type query string false "事件类型"
// @Param start_at query string false "开始时间 RFC3339"
// @Param end_at query string false "结束时间 RFC3339"
// @Success 200 {object} v1.ApiResponse[eventV1.ListDeviceEventsResponse]
// @Router /device-events [get]
func (h *DeviceEventHandler) ListDeviceEvents(c *gin.Context) {
	pageNumber, pageSize := page(c, 20)
	startAt, err := eventTime(c.Query("start_at"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "invalid start_at: " + err.Error()})
		return
	}
	endAt, err := eventTime(c.Query("end_at"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "invalid end_at: " + err.Error()})
		return
	}
	events, total, err := h.svc.List(c, pageNumber, pageSize, c.Query("keyword"), c.Query("device_key"), c.Query("event_type"), startAt, endAt)
	if err != nil { deviceError(c, err); return }
	out := make([]eventV1.DeviceEvent, 0, len(events))
	for _, event := range events { out = append(out, deviceEventResponse(event)) }
	v1.HandleSuccess(c, eventV1.ListDeviceEventsResponse{Events: out, Total: total, Page: pageNumber, PageSize: pageSize})
}

func eventTime(value string) (*time.Time, error) {
	if value == "" { return nil, nil }
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil { return nil, err }
	return &parsed, nil
}
