package handler

import (
	"aiot-backend/api/v1"
	"aiot-backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

// Register godoc
// @Summary Register user
// @Schemes
// @Description Registers user.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "params"
// @Success 200 {object} v1.Response "Successful response"
// @Router /register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	req := new(v1.RegisterRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Register(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// Login godoc
// @Summary Log in
// @Schemes
// @Description Authenticates a user and returns an access token.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "params"
// @Success 200 {object} v1.ApiResponse[v1.LoginResponseData] "Successful response"
// @Router /login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.LoginResponseData{
		AccessToken: token,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Schemes
// @Description Returns the profile of the authenticated user.
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.ApiResponse[v1.GetProfileResponseData] "Successful response"
// @Router /user [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	user, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Schemes
// @Description Updates the profile of the authenticated user.
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateProfileRequest true "params"
// @Success 200 {object} v1.Response "Successful response"
// @Router /user [put]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)

	var req v1.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.UpdateProfile(ctx, userId, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// ListMyOrganizations godoc
// @Summary List user organizations
// @Schemes
// @Description Lists organizations available to the authenticated user.
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.ApiResponse[[]v1.OrganizationItem] "Successful response"
// @Router /organizations [get]
func (h *UserHandler) ListMyOrganizations(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	items, err := h.userService.ListMyOrganizations(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, items)
}

// SwitchOrganization godoc
// @Summary Switch organization
// @Schemes
// @Description Switches the current organization and returns a new access token.
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.SwitchOrgRequest true "params"
// @Success 200 {object} v1.ApiResponse[v1.SwitchOrgResponseData] "Successful response"
// @Router /auth/switch-org [post]
func (h *UserHandler) SwitchOrganization(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	var req v1.SwitchOrgRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	token, err := h.userService.SwitchOrganization(ctx, userId, req.OrgId)
	if err != nil {
		if err == v1.ErrForbidden {
			v1.HandleError(ctx, http.StatusForbidden, v1.ErrForbidden, nil)
			return
		}
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, v1.SwitchOrgResponseData{
		AccessToken: token,
	})
}
