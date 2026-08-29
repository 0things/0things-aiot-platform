package router

import "github.com/gin-gonic/gin"

func InitProductTSLRouter(deps RouterDeps, r *gin.RouterGroup) {
	tsl := deps.ProductTSLHandler
	r.POST("/products/:productKey/tsl", tsl.Put)
	r.GET("/products/:productKey/tsl", tsl.Get)
	r.PUT("/products/:productKey/tsl", tsl.Put)
	r.DELETE("/products/:productKey/tsl", tsl.Delete)
}
