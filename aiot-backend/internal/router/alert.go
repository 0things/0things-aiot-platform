package router

import "github.com/gin-gonic/gin"

func InitAlertRouter(deps RouterDeps, r *gin.RouterGroup) {
	alert := deps.AlertHandler
	r.GET("/alerts", alert.ListAlerts)
	r.GET("/alerts/:id", alert.GetAlert)
	r.POST("/alerts/:id/ack", func(c *gin.Context) {
		alert.AlertStatus(c, "acknowledged")
	})
	r.POST("/alerts/:id/resolve", func(c *gin.Context) {
		alert.AlertStatus(c, "resolved")
	})
}
