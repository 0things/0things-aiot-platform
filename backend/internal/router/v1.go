package router

import (
	"aiot-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// InitV1Routers registers each independently owned v1 route group.
func InitV1Routers(deps RouterDeps, r *gin.RouterGroup) {
	InitOrganizationRouter(deps, r)
	protected := r.Group("/")
	protected.Use(middleware.StrictAuth(deps.Logto, deps.Logger))
	InitCategoryRouter(deps, protected)
	InitProductRouter(deps, protected)
	InitProductTSLRouter(deps, protected)
	InitProductMessageParserRouter(deps, protected)
	InitDeviceRouter(deps, protected)
	InitDeviceGroupRouter(deps, protected)
	InitOTARouter(deps, protected)
	InitFileRouter(deps, protected)
	InitDeviceEventRouter(deps, protected)
	InitTelemetryRouter(deps, protected)
	InitRuleNodeDefinitionRouter(deps, protected)
	InitRuleChainRouter(deps, protected)
}
