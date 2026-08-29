package router

import "github.com/gin-gonic/gin"

func InitCategoryRouter(deps RouterDeps, r *gin.RouterGroup) {
	r.GET("/categories/tree", deps.CategoryHandler.Tree)
}
