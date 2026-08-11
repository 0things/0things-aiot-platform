package router

import "github.com/gin-gonic/gin"

func InitProductTSLRouter(deps RouterDeps, r *gin.RouterGroup) {
	tsl := deps.ProductTSLHandler
	r.POST("/products/:id/tsl", tsl.Put)
	r.GET("/products/:id/tsl", tsl.Get)
	r.PUT("/products/:id/tsl", tsl.Put)
	r.DELETE("/products/:id/tsl", tsl.Delete)
}
