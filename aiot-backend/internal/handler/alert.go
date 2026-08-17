package handler

import (
	"fmt"

	alertV1 "0things-backend/api/alert/v1"
	v1 "0things-backend/api/v1"
	"0things-backend/internal/model"
	"0things-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	*Handler
	svc service.AlertServiceInterface
}

func NewAlertHandler(h *Handler, svc service.AlertServiceInterface) *AlertHandler {
	return &AlertHandler{Handler: h, svc: svc}
}

func alertResponse(alert model.Alert) alertV1.Alert {
	return alertV1.Alert{ID: alert.ID, RuleID: alert.RuleID, RuleName: alert.RuleName, DeviceKey: alert.DeviceKey, Severity: alert.Severity, Status: alert.Status, Summary: alert.Summary, Payload: fmt.Sprint(raw(alert.Payload)), Count: alert.Count, RaisedAt: alert.RaisedAt, LastRaisedAt: alert.LastRaisedAt, AckAt: alert.AckAt, AckBy: alert.AckBy, ResolvedAt: alert.ResolvedAt, ResolvedBy: alert.ResolvedBy, CreatedAt: alert.CreatedAt, UpdatedAt: alert.UpdatedAt}
}
// ListAlerts godoc
// @Summary 获取告警列表
// @Schemes
// @Description 分页获取告警列表，支持按 status、severity、device_key 过滤
// @Tags 告警模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param status query string false "状态"
// @Param severity query string false "严重程度"
// @Param device_key query string false "设备 Key"
// @Success 200 {object} v1.ApiResponse[alertV1.ListAlertsResponse]
// @Router /alerts [get]
func (h *AlertHandler) ListAlerts(c *gin.Context) {
	pageNumber, pageSize := page(c, 20)
	alerts, total, err := h.svc.List(c, pageNumber, pageSize, c.Query("status"), c.Query("severity"), c.Query("device_key"))
	if err != nil {
		deviceError(c, err)
		return
	}
	out := make([]alertV1.Alert, 0, len(alerts))
	for _, alert := range alerts {
		out = append(out, alertResponse(alert))
	}
	v1.HandleSuccess(c, alertV1.ListAlertsResponse{Alerts: out, Total: total, Page: pageNumber, PageSize: pageSize})
}

// GetAlert godoc
// @Summary 获取告警详情
// @Schemes
// @Description 通过告警 ID 获取告警详情
// @Tags 告警模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "告警 ID"
// @Success 200 {object} v1.ApiResponse[alertV1.GetAlertResponse]
// @Router /alerts/{id} [get]
func (h *AlertHandler) GetAlert(c *gin.Context) {
	alertID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	alert, err := h.svc.Get(c, alertID)
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, alertV1.GetAlertResponse{Alert: alertResponse(*alert)})
}

// AckAlert godoc
// @Summary 确认告警
// @Schemes
// @Description 通过告警 ID 确认（acknowledged）该告警
// @Tags 告警模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "告警 ID"
// @Success 200 {object} v1.ApiResponse[alertV1.GetAlertResponse]
// @Router /alerts/{id}/ack [post]
func (h *AlertHandler) AckAlert(c *gin.Context) {
	h.AlertStatus(c, "acknowledged")
}

// ResolveAlert godoc
// @Summary 处理告警
// @Schemes
// @Description 通过告警 ID 处理（resolved）该告警
// @Tags 告警模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "告警 ID"
// @Success 200 {object} v1.ApiResponse[alertV1.GetAlertResponse]
// @Router /alerts/{id}/resolve [post]
func (h *AlertHandler) ResolveAlert(c *gin.Context) {
	h.AlertStatus(c, "resolved")
}

func (h *AlertHandler) AlertStatus(c *gin.Context, status string) {
	alertID, err := id(c)
	if err != nil {
		deviceError(c, err)
		return
	}
	alert, err := h.svc.SetStatus(c, alertID, status)
	if err != nil {
		deviceError(c, err)
		return
	}
	v1.HandleSuccess(c, alertV1.GetAlertResponse{Alert: alertResponse(*alert)})
}
