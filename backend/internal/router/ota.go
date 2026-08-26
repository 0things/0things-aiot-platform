package router

import "github.com/gin-gonic/gin"

func InitOTARouter(deps RouterDeps, r *gin.RouterGroup) {
	ota := deps.OTAHandler
	r.POST("/ota-packages", ota.CreateOTA)
	r.GET("/ota-packages", ota.ListOTA)
	r.GET("/ota-packages/:uuid", ota.GetOTA)
	r.PUT("/ota-packages/:uuid", ota.UpdateOTA)
	r.DELETE("/ota-packages/:uuid", ota.DeleteOTA)
	r.POST("/ota-packages/:uuid/batch-upgrade", ota.BatchUpgradeOTA)
	r.POST("/ota-packages/:uuid/batches/:batchId/cancel", ota.CancelBatch)
	r.POST("/ota-packages/:uuid/batches/:batchId/retry", ota.RetryBatch)
	r.POST("/ota-packages/:uuid/report", ota.ReportOTAStatus)
	r.GET("/ota-packages/:uuid/upgrade-statistics", ota.OTAStats)
	r.GET("/ota-packages/:uuid/batches", ota.OTABatches)
	r.GET("/ota-packages/:uuid/device-deployments", ota.OTADeployments)
}
