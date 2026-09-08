package router

import (
	"net/http"

	"transport-http/internal/handler"
	"transport-http/internal/middleware"
	"transport-http/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Router struct {
	deviceHandler *handler.DeviceHandler
	config        *viper.Viper
	logger        *log.Logger
}

func NewRouter(deviceHandler *handler.DeviceHandler, conf *viper.Viper, logger *log.Logger) *Router {
	return &Router{
		deviceHandler: deviceHandler,
		config:        conf,
		logger:        logger,
	}
}

func (r *Router) Register(e *gin.Engine) {
	e.Use(gin.Recovery())

	// Health check
	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "transport-http"})
	})

	authMiddleware := middleware.DeviceAuthMiddleware(r.config, r.logger)

	// API v1 device routes
	apiV1 := e.Group("/api/v1")
	apiV1.Use(authMiddleware)
	{
		apiV1.POST("/:deviceKey/telemetry", r.deviceHandler.PostTelemetry)
		apiV1.POST("/:deviceKey/attributes", r.deviceHandler.PostAttributes)
		apiV1.POST("/:deviceKey/events/:eventType", r.deviceHandler.PostEvent)
		apiV1.POST("/:deviceKey/ota/progress", r.deviceHandler.PostOtaProgress)
	}

	// Legacy device ingress route
	e.POST("/v1/device-ingress/:deviceKey", authMiddleware, r.deviceHandler.DeviceIngressLegacy)
}
