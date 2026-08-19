package router

import "github.com/gin-gonic/gin"

func InitOTARouter(deps RouterDeps, r *gin.RouterGroup) {
	ota := deps.OTAHandler
	r.POST("/ota-packages", ota.CreateOTA)
	r.GET("/ota-packages", ota.ListOTA)
	r.GET("/ota-packages/:id", ota.GetOTA)
	r.PUT("/ota-packages/:id", ota.UpdateOTA)
	r.DELETE("/ota-packages/:id", ota.DeleteOTA)
	r.POST("/ota-packages/:id/deploy", ota.DeployOTA)
	r.GET("/ota-packages/:id/upgrade-statistics", ota.OTAStats)
	r.GET("/ota-packages/:id/batches", ota.OTABatches)
	r.GET("/ota-packages/:id/device-deployments", ota.OTADeployments)
}
