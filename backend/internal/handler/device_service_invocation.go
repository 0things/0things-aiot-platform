package handler

import (
	"net/http"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/dto"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type DeviceServiceInvocationHandler struct {
	*Handler
	svc service.DeviceServiceInvocationServiceInterface
}

func NewDeviceServiceInvocationHandler(h *Handler, svc service.DeviceServiceInvocationServiceInterface) *DeviceServiceInvocationHandler {
	return &DeviceServiceInvocationHandler{Handler: h, svc: svc}
}

func deviceServiceInvocationJSON(item model.DeviceServiceInvocation) v1.DeviceServiceInvocation {
	return v1.DeviceServiceInvocation{UUID: item.UUID, InvokedAt: item.InvokedAt, ServiceIdentifier: item.ServiceIdentifier, ServiceName: item.ServiceName, InputParams: item.InputParams, OutputParams: item.OutputParams}
}

// List godoc
// @Summary List thing-model service invocation records
// @Description Lists paginated thing-model service invocation records for one device.
// @Tags Devices
// @Produce json
// @Security Bearer
// @Param deviceKey path string true "Device key"
// @Param request query v1.ListDeviceServiceInvocationsRequest false "Query parameters"
// @Success 200 {object} v1.ApiResponse[v1.ListDeviceServiceInvocationsResponse] "Successful response"
// @Router /devices/{deviceKey}/thing-model-service-invocations [get]
func (h *DeviceServiceInvocationHandler) List(c *gin.Context) {
	var req v1.ListDeviceServiceInvocationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	page, pageSize := pageRequest(req.PageRequest, 20)
	query := dto.ListDeviceServiceInvocationsQuery{
		DeviceKey:         c.Param("deviceKey"),
		ServiceIdentifier: req.ServiceIdentifier,
		StartAt:           req.StartAt,
		EndAt:             req.EndAt,
		Page:              page,
		PageSize:          pageSize,
	}
	items, total, err := h.svc.List(c, query)
	if err != nil {
		deviceError(c, err)
		return
	}
	out := make([]v1.DeviceServiceInvocation, len(items))
	for i, item := range items {
		out[i] = deviceServiceInvocationJSON(item)
	}
	v1.HandleSuccess(c, v1.ListDeviceServiceInvocationsResponse{Invocations: out, Total: total, Page: page, PageSize: pageSize})
}
