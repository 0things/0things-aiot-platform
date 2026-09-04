package handler

import (
	"net/http"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/dto"
	"aiot-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type DeviceEventHandler struct {
	*Handler
	svc service.DeviceEventServiceInterface
}

func NewDeviceEventHandler(h *Handler, svc service.DeviceEventServiceInterface) *DeviceEventHandler {
	return &DeviceEventHandler{Handler: h, svc: svc}
}

func deviceEventResponse(event dto.DeviceEventListItem) v1.DeviceEvent {
	return v1.DeviceEvent{
		ID:              event.ID,
		UUID:            event.UUID,
		DeviceKey:       event.DeviceKey,
		DeviceName:      event.DeviceName,
		EventIdentifier: event.EventIdentifier,
		EventType:       event.EventType,
		EventAt:         event.EventAt,
		Data:            string(event.Data),
	}
}

// ListDeviceEvents godoc
// @Summary List device events
// @Description Paginated query for thing model device events with keyword, deviceKey, eventType and date range filters
// @Tags Events
// @Produce json
// @Security Bearer
// @Param request query v1.ListDeviceEventsRequest false "Query parameters"
// @Success 200 {object} v1.ApiResponse[v1.ListDeviceEventsResponse] "Successful response"
// @Router /device-events [get]
func (h *DeviceEventHandler) ListDeviceEvents(c *gin.Context) {
	var req v1.ListDeviceEventsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	events, total, err := h.svc.List(c, &req)
	if err != nil {
		deviceError(c, err)
		return
	}
	out := make([]v1.DeviceEvent, 0, len(events))
	for _, event := range events {
		out = append(out, deviceEventResponse(event))
	}
	v1.HandleSuccess(c, v1.ListDeviceEventsResponse{Events: out, Total: total, Page: req.Page, PageSize: req.PageSize})
}
