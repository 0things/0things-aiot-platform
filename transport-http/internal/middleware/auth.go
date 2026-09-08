package middleware

import (
	"net/http"
	"strings"

	"transport-http/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// DeviceAuthMiddleware provides configurable authentication for device HTTP ingress.
// Supports X-Device-Token header or Authorization: Bearer <token>.
func DeviceAuthMiddleware(config *viper.Viper, logger *log.Logger) gin.HandlerFunc {
	authEnabled := config.GetBool("auth.enabled")
	expectedToken := config.GetString("auth.token")

	return func(c *gin.Context) {
		if !authEnabled || expectedToken == "" {
			c.Next()
			return
		}

		token := c.GetHeader("X-Device-Token")
		if token == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if token == "" || token != expectedToken {
			logger.Warn("unauthorized device HTTP ingress request",
				zap.String("device_key", c.Param("deviceKey")),
				zap.String("client_ip", c.ClientIP()),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorized: missing or invalid device token",
			})
			return
		}

		c.Next()
	}
}
