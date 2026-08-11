package router

import "github.com/gin-gonic/gin"

// InitV1Routers registers each independently owned v1 route group.
func InitV1Routers(deps RouterDeps, r *gin.RouterGroup) {
	InitProductRouter(deps, r)
	InitProductTSLRouter(deps, r)
	InitDeviceRouter(deps, r)
	InitRuleRouter(deps, r)
	InitAlertRouter(deps, r)
	InitOTARouter(deps, r)
}
