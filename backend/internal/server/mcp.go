package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/handler"
	"aiot-backend/internal/tenant"
	"aiot-backend/pkg/jwt"
	"aiot-backend/pkg/log"
	mcptransport "aiot-backend/pkg/server/mcp"
	nethttp "net/http"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const mcpToolTimeout = 10 * time.Second

func NewMCPServer(config *viper.Viper, logger *log.Logger, handler *handler.MCPHandler) *mcpserver.MCPServer {
	server := setupSrv(config, logger)
	registerTools(server, handler)
	return server
}

func NewMCPTransportServer(config *viper.Viper, logger *log.Logger, j *jwt.JWT, server *mcpserver.MCPServer) *mcptransport.Server {
	options := make([]mcptransport.Option, 0, 3)
	if config.GetBool("mcp.stdio_enabled") {
		options = append(options, mcptransport.WithStdioSrv())
	}
	if config.GetBool("mcp.sse_enabled") {
		options = append(options, mcptransport.WithSSESrv(config.GetString("mcp.sse_addr")))
	}
	if config.GetBool("mcp.streamable_http_enabled") {
		authMiddleware := func(next nethttp.Handler) nethttp.Handler {
			return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
				tokenString := r.Header.Get("Authorization")
				if tokenString == "" {
					logger.Warn("MCP HTTP request rejected: missing token", zap.String("path", r.URL.Path))
					nethttp.Error(w, "Unauthorized: missing Authorization token", nethttp.StatusUnauthorized)
					return
				}
				claims, err := j.ParseToken(tokenString)
				if err != nil || claims.OrganizationID <= 0 {
					logger.Warn("MCP HTTP request rejected: invalid token", zap.String("path", r.URL.Path), zap.Error(err))
					nethttp.Error(w, "Unauthorized: invalid or expired token", nethttp.StatusUnauthorized)
					return
				}
				ctx := tenant.WithTenant(r.Context(), claims.OrganizationID)
				ctx = context.WithValue(ctx, "claims", claims)
				ctx = context.WithValue(ctx, "user_id", claims.UserId)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}
		options = append(options, mcptransport.WithStreamableHTTPSrv(config.GetString("mcp.streamable_http_addr"), authMiddleware))
	}
	return mcptransport.NewServer(server, options...)
}

func setupSrv(config *viper.Viper, logger *log.Logger) *mcpserver.MCPServer {
	name := config.GetString("mcp.name")
	if name == "" {
		name = "0things-iot-mcp"
	}
	version := config.GetString("mcp.version")
	if version == "" {
		version = "0.1.0"
	}

	return mcpserver.NewMCPServer(
		name,
		version,
		mcpserver.WithToolCapabilities(false),
		mcpserver.WithRecovery(),
		mcpserver.WithHooks(newHooks(logger)),
		mcpserver.WithToolHandlerMiddleware(func(next mcpserver.ToolHandlerFunc) mcpserver.ToolHandlerFunc {
			return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				timeoutCtx, cancel := context.WithTimeout(ctx, mcpToolTimeout)
				defer cancel()
				return next(timeoutCtx, request)
			}
		}),
	)
}

func newHooks(logger *log.Logger) *mcpserver.Hooks {
	hooks := &mcpserver.Hooks{}
	startedAt := sync.Map{}

	hooks.AddBeforeAny(func(_ context.Context, id any, method mcp.MCPMethod, _ any) {
		startedAt.Store(fmt.Sprint(id), time.Now())
		logger.Debug("MCP request started", zap.String("method", string(method)), zap.Any("request_id", id))
	})
	hooks.AddOnSuccess(func(_ context.Context, id any, method mcp.MCPMethod, _ any, _ any) {
		logger.Info("MCP request completed",
			zap.String("method", string(method)),
			zap.Any("request_id", id),
			zap.Duration("duration", hookDuration(&startedAt, id)),
		)
	})
	hooks.AddOnError(func(_ context.Context, id any, method mcp.MCPMethod, _ any, err error) {
		logger.Error("MCP request failed",
			zap.String("method", string(method)),
			zap.Any("request_id", id),
			zap.Duration("duration", hookDuration(&startedAt, id)),
			zap.Error(err),
		)
	})
	hooks.AddBeforeCallTool(func(ctx context.Context, id any, request *mcp.CallToolRequest) {
		orgID := tenant.GetOrganizationID(ctx)
		var userID string
		if claims, ok := ctx.Value("claims").(*jwt.MyCustomClaims); ok && claims != nil {
			userID = claims.UserId
		}
		logger.Info("MCP tool call started",
			zap.Any("request_id", id),
			zap.Int64("organization_id", orgID),
			zap.String("user_id", userID),
			zap.String("tool", request.Params.Name),
		)
	})
	hooks.AddAfterCallTool(func(ctx context.Context, id any, request *mcp.CallToolRequest, _ any) {
		orgID := tenant.GetOrganizationID(ctx)
		logger.Info("MCP tool call completed",
			zap.Any("request_id", id),
			zap.Int64("organization_id", orgID),
			zap.String("tool", request.Params.Name),
			zap.Duration("duration", hookDuration(&startedAt, id)),
		)
	})
	return hooks
}

func hookDuration(startedAt *sync.Map, id any) time.Duration {
	key := fmt.Sprint(id)
	started, ok := startedAt.LoadAndDelete(key)
	if !ok {
		return 0
	}
	return time.Since(started.(time.Time))
}

func registerTools(server *mcpserver.MCPServer, handler *handler.MCPHandler) {
	server.AddTool(
		mcp.NewTool("iot_query_devices",
			mcp.WithDescription("Query devices in the current organization."),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithNumber("productId", mcp.Description("Optional product ID.")),
			mcp.WithString("status", mcp.Description("Optional device status.")),
			mcp.WithString("keyword", mcp.Description("Optional device name or key keyword.")),
			mcp.WithNumber("limit", mcp.Description("Maximum number of devices, from 1 to 100.")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args v1.QueryDevicesRequest
			if err := request.BindArguments(&args); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			result, err := handler.QueryDevices(ctx, args)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(result)
		},
	)

	server.AddTool(
		mcp.NewTool("iot_get_device_detail",
			mcp.WithDescription("Get a device with its tags, shadow, and thing-model properties."),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithString("deviceKey", mcp.Required(), mcp.Description("Target device key.")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args v1.GetDeviceDetailRequest
			if err := request.BindArguments(&args); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			result, err := handler.GetDeviceDetail(ctx, args)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(result)
		},
	)

	server.AddTool(
		mcp.NewTool("iot_query_telemetry_history",
			mcp.WithDescription("Query telemetry points and a numeric summary for one device property."),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithString("deviceKey", mcp.Required(), mcp.Description("Target device key.")),
			mcp.WithString("property", mcp.Required(), mcp.Description("Thing-model property identifier.")),
			mcp.WithNumber("startTime", mcp.Description("Optional inclusive Unix millisecond timestamp.")),
			mcp.WithNumber("endTime", mcp.Description("Optional inclusive Unix millisecond timestamp.")),
			mcp.WithNumber("limit", mcp.Description("Maximum points, from 1 to 100.")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args v1.QueryTelemetryHistoryRequest
			if err := request.BindArguments(&args); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			result, err := handler.QueryTelemetryHistory(ctx, args)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultJSON(result)
		},
	)
}
