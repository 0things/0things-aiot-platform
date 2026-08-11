package handler

import (
	"fmt"

	alertV1 "0things-backend/api/alert/v1"
	"0things-backend/internal/model"
	"0things-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	*Handler
	svc *service.AlertService
}

func NewAlertHandler(h *Handler, svc *service.AlertService) *AlertHandler {
	return &AlertHandler{Handler: h, svc: svc}
}

func alertResponse(alert model.Alert) alertV1.Alert {
	return alertV1.Alert{ID: alert.ID, RuleID: alert.RuleID, RuleName: alert.RuleName, DeviceKey: alert.DeviceKey, Severity: alert.Severity, Status: alert.Status, Summary: alert.Summary, Payload: fmt.Sprint(raw(alert.Payload)), Count: alert.Count, RaisedAt: alert.RaisedAt, LastRaisedAt: alert.LastRaisedAt, AckAt: alert.AckAt, AckBy: alert.AckBy, ResolvedAt: alert.ResolvedAt, ResolvedBy: alert.ResolvedBy, CreatedAt: alert.CreatedAt, UpdatedAt: alert.UpdatedAt}
}
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
	c.JSON(200, alertV1.ListAlertsResponse{Alerts: out, Total: total, Page: pageNumber, PageSize: pageSize})
}
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
	c.JSON(200, alertV1.GetAlertResponse{Alert: alertResponse(*alert)})
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
	c.JSON(200, alertV1.GetAlertResponse{Alert: alertResponse(*alert)})
}
