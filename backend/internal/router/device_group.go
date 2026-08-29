package router

import "github.com/gin-gonic/gin"

func InitDeviceGroupRouter(deps RouterDeps, r *gin.RouterGroup) {
	group := deps.DeviceGroupHandler
	r.POST("/device-groups", group.Create)
	r.GET("/device-groups", group.List)
	r.POST("/device-groups/preview", group.Preview)
	r.GET("/device-groups/:groupUuid", group.Get)
	r.PUT("/device-groups/:groupUuid", group.Update)
	r.DELETE("/device-groups/:groupUuid", group.Delete)
	r.GET("/device-groups/:groupUuid/devices", group.ListDevices)
	r.POST("/device-groups/:groupUuid/devices", group.AddDevices)
	r.DELETE("/device-groups/:groupUuid/devices", group.RemoveDevices)
	r.POST("/device-groups/:groupUuid/preview", group.PreviewSaved)
}
