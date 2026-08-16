package router

import (
	"github.com/gin-gonic/gin"
)

func InitSceneLinkageRouter(deps RouterDeps, r *gin.RouterGroup) {
	sl := deps.SceneLinkageHandler
	r.GET("/scene-linkages", sl.ListSceneLinkages)
	r.POST("/scene-linkages", sl.CreateSceneLinkage)
	r.GET("/scene-linkages/:id", sl.GetSceneLinkage)
	r.PUT("/scene-linkages/:id", sl.UpdateSceneLinkage)
	r.DELETE("/scene-linkages/:id", sl.DeleteSceneLinkage)

	detail := deps.SceneLinkageDetailHandler
	r.GET("/scene-linkages/:id/detail", detail.GetSceneLinkageDetail)
	r.POST("/scene-linkages/:id/detail", detail.CreateSceneLinkageDetail)
	r.PUT("/scene-linkages/:id/detail", detail.UpdateSceneLinkageDetail)
}
