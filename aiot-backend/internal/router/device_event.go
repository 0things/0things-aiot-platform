package router

import "github.com/gin-gonic/gin"

func InitDeviceEventRouter(deps RouterDeps, r *gin.RouterGroup) {
	r.GET("/device-events", deps.DeviceEventHandler.ListDeviceEvents)
}
