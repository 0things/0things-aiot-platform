package router

import (
	"aiot-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitOrganizationRouter(deps RouterDeps, r *gin.RouterGroup) {
	identity := r.Group("/")
	identity.Use(middleware.IdentityAuth(deps.Logto, deps.Logger))
	identity.POST("/me/organization", deps.OrganizationHandler.EnsureOrganization)
}
