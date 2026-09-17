package middleware

import (
	"net/http"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/tenant"
	"aiot-backend/pkg/log"
	"aiot-backend/pkg/logto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ClaimsKey = "claims"
	UserIDKey = "user_id"
)

func StrictAuth(verifier *logto.Verifier, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !authenticate(ctx, verifier, logger) {
			return
		}
		claimsVal, exists := ctx.Get(ClaimsKey)
		if !exists {
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}
		claims, ok := claimsVal.(*logto.Claims)
		if !ok || claims == nil {
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}
		organizationID, err := verifier.OrganizationID(claims)
		if err != nil {
			v1.HandleError(ctx, http.StatusForbidden, v1.ErrForbidden, nil)
			ctx.Abort()
			return
		}
		ctx.Set(string(tenant.OrganizationKey), organizationID)
		ctx.Request = ctx.Request.WithContext(tenant.WithOrganization(ctx.Request.Context(), organizationID))
		ctx.Next()
	}
}

// IdentityAuth verifies a user token without requiring an organization claim.
// It is only used by the organization bootstrap endpoint.
func IdentityAuth(verifier *logto.Verifier, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !authenticate(ctx, verifier, logger) {
			return
		}
		ctx.Next()
	}
}

func authenticate(ctx *gin.Context, verifier *logto.Verifier, logger *log.Logger) bool {
	tokenString := ctx.Request.Header.Get("Authorization")
	if tokenString == "" {
		logger.WithContext(ctx).Warn("No token", zap.Any("data", map[string]interface{}{
			"url":    ctx.Request.URL,
			"params": ctx.Params,
		}))
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		ctx.Abort()
		return false
	}

	claims, err := verifier.Parse(ctx.Request.Context(), tokenString)
	if err != nil {
		logger.WithContext(ctx).Error("token error", zap.Any("data", map[string]interface{}{
			"url":    ctx.Request.URL,
			"params": ctx.Params,
		}), zap.Error(err))
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		ctx.Abort()
		return false
	}

	ctx.Set(ClaimsKey, claims)
	ctx.Set(UserIDKey, claims.Subject)
	recoveryLoggerFunc(ctx, logger)
	return true
}

func recoveryLoggerFunc(ctx *gin.Context, logger *log.Logger) {
	if val, exists := ctx.Get(ClaimsKey); exists {
		if userInfo, ok := val.(*logto.Claims); ok && userInfo != nil {
			logger.WithValue(ctx, zap.String("UserId", userInfo.Subject))
		}
	}
}
