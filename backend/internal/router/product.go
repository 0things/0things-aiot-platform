package router

import "github.com/gin-gonic/gin"

func InitProductRouter(deps RouterDeps, r *gin.RouterGroup) {
	product := deps.ProductHandler
	r.POST("/products", product.Create)
	r.GET("/products", product.List)
	r.GET("/products/:id", product.Get)
	r.GET("/products/key/:productKey", product.GetByKey)
	r.PUT("/products/key/:productKey", product.Update)
	r.DELETE("/products/:id", product.Delete)
	r.POST("/products/:id/restore", product.Restore)
}
