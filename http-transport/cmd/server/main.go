package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"http-transport/internal/handler"
	"http-transport/internal/kafka"
	"http-transport/pkg/config"
	"http-transport/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// main 是 0things HTTP 传输微服务的启动入口。
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)

	logger.Info("starting 0things HTTP Transport Service...")

	// 1. 初始化 Kafka Producer
	producer, cleanupProducer, err := kafka.NewProducer(conf, logger.Logger)
	if err != nil {
		logger.Fatal("failed to initialize kafka producer", zap.Error(err))
	}
	defer cleanupProducer()

	// 2. 初始化 Gin Web 引擎
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	deviceHandler := handler.NewDeviceHandler(producer, logger.Logger)

	// 注册标准设备上报路由
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/:deviceKey/telemetry", deviceHandler.PostTelemetry)
		apiV1.POST("/:deviceKey/attributes", deviceHandler.PostAttributes)
		apiV1.POST("/:deviceKey/events/:eventType", deviceHandler.PostEvent)
		apiV1.POST("/:deviceKey/ota/progress", deviceHandler.PostOtaProgress) // 新增 OTA 进度上报接口
	}
	// 兼容老网关路由
	r.POST("/v1/device-ingress/:deviceKey", deviceHandler.DeviceIngressLegacy)

	// 健康检查探针接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "http-transport"})
	})

	host := conf.GetString("http.host")
	port := conf.GetInt("http.port")
	if port == 0 {
		port = 8081
	}
	addr := fmt.Sprintf("%s:%d", host, port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("HTTP Transport listening on", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server listen failed", zap.Error(err))
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	logger.Info("HTTP Transport Service stopped gracefully")
}
