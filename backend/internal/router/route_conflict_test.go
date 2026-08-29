package router

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestDeviceSubresourceRoutesDoNotConflictWithResourceRoutes(t *testing.T) {
	r := gin.New()
	group := r.Group("/v1")
	InitProductRouter(RouterDeps{}, group)
	InitProductTSLRouter(RouterDeps{}, group)
	InitProductMessageParserRouter(RouterDeps{}, group)
	InitDeviceRouter(RouterDeps{}, group)
}
