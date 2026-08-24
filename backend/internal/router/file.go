package router

import "github.com/gin-gonic/gin"

func InitFileRouter(deps RouterDeps, r *gin.RouterGroup) {
	r.POST("/files/ota", deps.FileHandler.UploadOTAFile)
}
