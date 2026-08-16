package router

import (
	"github.com/gin-gonic/gin"
)

// InitV1Routers registers each independently owned v1 route group.
func InitV1Routers(deps RouterDeps, r *gin.RouterGroup) {
	protected := r.Group("/")
	// protected.Use(middleware.StrictAuth(deps.JWT, deps.Logger))
	InitProductRouter(deps, protected)
	InitProductTSLRouter(deps, protected)
	InitDeviceRouter(deps, protected)
	InitSceneLinkageRouter(deps, protected)
	InitAlertRouter(deps, protected)
	InitOTARouter(deps, protected)
	InitDeviceEventRouter(deps, protected)
}
