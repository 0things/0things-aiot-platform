package router

import "github.com/gin-gonic/gin"

func InitProductRouter(deps RouterDeps, r *gin.RouterGroup) {
	product := deps.ProductHandler
	r.POST("/products", product.Create)
	r.GET("/products", product.List)
	r.GET("/products/options", product.Options)
	r.GET("/products/:productKey", product.Get)
	r.PUT("/products/:productKey", product.Update)
	r.DELETE("/products/:productKey", product.Delete)
}
