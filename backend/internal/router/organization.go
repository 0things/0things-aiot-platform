package router

import (
	"aiot-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitOrganizationRouter(deps RouterDeps, r *gin.RouterGroup) {
	identity := r.Group("/")
	identity.Use(middleware.IdentityAuth(deps.Logto, deps.Logger))
	identity.POST("/organizations", deps.OrganizationHandler.CreateOrganization)
	identity.GET("/organizations/invitations/:id", deps.OrganizationHandler.GetInvitation)
	identity.POST("/organizations/invitations/:id/accept", deps.OrganizationHandler.AcceptInvitation)
}

func InitOrganizationProtectedRouter(deps RouterDeps, r *gin.RouterGroup) {
	org := deps.OrganizationHandler
	r.POST("/organizations/invitations", org.CreateInvitation)
	r.GET("/organizations/invitations", org.ListInvitations)
	r.DELETE("/organizations/invitations/:id", org.DeleteInvitation)
	r.GET("/organizations/members", org.ListMembers)
	r.DELETE("/organizations/members/:userId", org.RemoveMember)
}
