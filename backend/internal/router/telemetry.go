package router

import (
	"github.com/gin-gonic/gin"
)

func InitTelemetryRouter(deps RouterDeps, r *gin.RouterGroup) {
	devices := r.Group("/devices/:deviceKey")
	{
		devices.GET("/telemetry/history", deps.TelemetryHandler.QueryHistory)
	}
}
