package router

import "github.com/gin-gonic/gin"

func InitAlertRouter(deps RouterDeps, r *gin.RouterGroup) {
	alert := deps.AlertHandler
	r.GET("/alerts", alert.ListAlerts)
	r.GET("/alerts/:id", alert.GetAlert)
	r.POST("/alerts/:id/ack", alert.AckAlert)
	r.POST("/alerts/:id/resolve", alert.ResolveAlert)
}
