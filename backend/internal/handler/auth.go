package handler

import (
	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/tenant"

	"github.com/gin-gonic/gin"
)

// GetAuthMe godoc
// @Summary Get the current authenticated user and organization
// @Tags Auth
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.ApiResponse[v1.AuthMeResp]
// @Router /auth/me [get]
func GetAuthMe(c *gin.Context) {
	organizationID := c.GetString(string(tenant.OrganizationKey))
	v1.HandleSuccess(c, v1.AuthMeResp{
		Subject:        userIDFromContext(c),
		OrganizationID: organizationID,
	})
}
