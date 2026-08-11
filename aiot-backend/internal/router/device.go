package router

import "github.com/gin-gonic/gin"

func InitDeviceRouter(deps RouterDeps, r *gin.RouterGroup) {
	device := deps.DeviceHandler
	r.POST("/devices", device.CreateDevice)
	r.GET("/devices", device.ListDevices)
	r.GET("/devices/:id", device.GetDevice)
	r.GET("/devices/key/:deviceKey", device.GetDeviceByKey)
	r.PUT("/devices/:id", device.UpdateDevice)
	r.DELETE("/devices/:id", device.DeleteDevice)
	r.POST("/devices/:id/activate", device.Activate)
	r.POST("/devices/:id/enabled", device.Enabled)
	r.POST("/devices/:id/restore", device.Restore)
	r.GET("/device-statistics", device.Stats)
	r.GET("/devices/:id/telemetry", device.Telemetry)
	r.GET("/devices/:id/mqtt-parameters", device.MQTT)
	r.GET("/devices/:id/tags", device.GetTags)
	r.PUT("/devices/:id/tags", device.PutTags)
	r.POST("/devices/:id/tags", device.PostTags)
	r.DELETE("/devices/:id/tags", device.DeleteTags)
	r.GET("/devices/:id/shadow", device.GetShadow)
	r.PUT("/devices/:id/shadow/desired", device.Desired)
	r.PUT("/devices/:id/shadow/reported", device.Reported)
	r.DELETE("/devices/:id/shadow/desired", device.ClearDesired)
	r.GET("/devices/:id/shadow/history", device.History)
	r.POST("/devices/:id/simulate-push", device.SimulatePush)
	r.GET("/devices/:id/push-records", device.PushRecords)
	r.GET("/devices/push-records/:pushRecordId", device.PushRecord)
	r.DELETE("/devices/:id/push-records", device.ClearPushRecords)
	r.POST("/devices/batch/upload", device.BatchUpload)
	r.GET("/devices/batch/template", device.BatchTemplate)
	r.POST("/devices/mock-kafka", device.MockKafka)
}
