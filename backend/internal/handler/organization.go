package handler

import (
	"net/http"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	*Handler
	svc service.OrganizationProvisioningServiceInterface
}

func NewOrganizationHandler(h *Handler, svc service.OrganizationProvisioningServiceInterface) *OrganizationHandler {
	return &OrganizationHandler{Handler: h, svc: svc}
}

// EnsureOrganization godoc
// @Summary Ensure the current user has a Logto organization
// @Tags Organizations
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.ApiResponse[map[string]string]
// @Router /me/organization [post]
func (h *OrganizationHandler) EnsureOrganization(c *gin.Context) {
	userID := userIDFromContext(c)
	if userID == "" {
		v1.HandleError(c, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	organizationID, err := h.svc.EnsureOrganization(c.Request.Context(), userID)
	if err != nil {
		v1.HandleError(c, http.StatusBadGateway, err, nil)
		return
	}
	v1.HandleSuccess(c, map[string]string{"organizationId": organizationID})
}
