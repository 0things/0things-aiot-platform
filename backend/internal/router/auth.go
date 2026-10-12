package router

import (
	"aiot-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func InitAuthRouter(r *gin.RouterGroup) {
	r.GET("/auth/me", handler.GetAuthMe)
}
