package handler

import (
	"aiot-backend/internal/middleware"
	"aiot-backend/pkg/log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{logger: logger}
}

// userIDFromContext reads the subject populated by the Logto middleware.
func userIDFromContext(ctx *gin.Context) string {
	value, exists := ctx.Get(middleware.UserIDKey)
	if !exists {
		return ""
	}
	userID, _ := value.(string)
	return userID
}
