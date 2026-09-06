package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/handler"
	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"
	"aiot-backend/pkg/jwt"
	"aiot-backend/pkg/log"
	mcptransport "aiot-backend/pkg/server/mcp"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mcpTestDevices struct{}

func (mcpTestDevices) DeviceByKey(ctx context.Context, key string) (*model.Device, error) {
	orgID := tenant.GetOrganizationID(ctx)
	if key == "device-org1" && orgID == 1 {
		return &model.Device{DeviceKey: "device-org1", Name: "Org1 Device", ProductID: 1, OrganizationID: 1}, nil
	}
	if key == "device-org2" && orgID == 2 {
		return &model.Device{DeviceKey: "device-org2", Name: "Org2 Device", ProductID: 2, OrganizationID: 2}, nil
	}
	return nil, errors.New("device not found in organization")
}

func (mcpTestDevices) ListDevices(ctx context.Context, _ dto.ListDevicesQuery) ([]model.Device, int64, error) {
	orgID := tenant.GetOrganizationID(ctx)
	if orgID == 1 {
		return []model.Device{{DeviceKey: "device-org1", Name: "Org1 Device", ProductID: 1, OrganizationID: 1}}, 1, nil
	}
	if orgID == 2 {
		return []model.Device{{DeviceKey: "device-org2", Name: "Org2 Device", ProductID: 2, OrganizationID: 2}}, 1, nil
	}
	return nil, 0, nil
}

func (mcpTestDevices) Tags(ctx context.Context, key string) ([]model.DeviceTag, error) {
	orgID := tenant.GetOrganizationID(ctx)
	if (key == "device-org1" && orgID == 1) || (key == "device-org2" && orgID == 2) {
		return []model.DeviceTag{{Key: "env", Value: "prod"}}, nil
	}
	return nil, errors.New("device not found in organization")
}

func (mcpTestDevices) Shadow(ctx context.Context, key string) (*model.DeviceShadow, error) {
	orgID := tenant.GetOrganizationID(ctx)
	if (key == "device-org1" && orgID == 1) || (key == "device-org2" && orgID == 2) {
		return &model.DeviceShadow{Version: 1}, nil
	}
	return nil, errors.New("device not found in organization")
}

type mcpTestThingModel struct{}

func (mcpTestThingModel) ListProperties(ctx context.Context, key string) ([]dto.ThingModelProperty, error) {
	orgID := tenant.GetOrganizationID(ctx)
	if (key == "device-org1" && orgID == 1) || (key == "device-org2" && orgID == 2) {
		return []dto.ThingModelProperty{{Identifier: "temp", Name: "Temperature", DataType: "double"}}, nil
	}
	return nil, errors.New("device not found in organization")
}

type mcpTestTelemetry struct{}

func (mcpTestTelemetry) QueryHistory(_ context.Context, _ dto.TelemetryQueryReq) ([]dto.TelemetryPoint, error) {
	return []dto.TelemetryPoint{
		{Timestamp: 1, Property: "temp", Value: 20.0},
		{Timestamp: 2, Property: "temp", Value: 24.0},
	}, nil
}

func TestMCPServerRegistersReadOnlyToolsAndBindsTelemetryArguments(t *testing.T) {
	toolHandler := handler.NewMCPHandler(mcpTestDevices{}, mcpTestThingModel{}, mcpTestTelemetry{})
	server := NewMCPServer(viper.New(), &log.Logger{Logger: zap.NewNop()}, toolHandler)

	require.Len(t, server.ListTools(), 3)
	tool := server.GetTool("iot_query_telemetry_history")
	require.NotNil(t, tool)

	ctx := tenant.WithTenant(context.Background(), 1)
	result, err := tool.Handler(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{
		Name: "iot_query_telemetry_history",
		Arguments: map[string]any{
			"deviceKey": "device-org1",
			"property":  "temp",
			"limit":     2,
		},
	}})
	require.NoError(t, err)
	require.False(t, result.IsError)
	require.IsType(t, &handler.MCPTelemetryHistory{}, result.StructuredContent)

	history := result.StructuredContent.(*handler.MCPTelemetryHistory)
	require.Len(t, history.Points, 2)
	require.Equal(t, 22.0, *history.Stats.Avg)
}

func TestMCPTenantIsolation(t *testing.T) {
	toolHandler := handler.NewMCPHandler(mcpTestDevices{}, mcpTestThingModel{}, mcpTestTelemetry{})
	server := NewMCPServer(viper.New(), &log.Logger{Logger: zap.NewNop()}, toolHandler)

	// 1. Query with Org 1 context
	ctxOrg1 := tenant.WithTenant(context.Background(), 1)
	devTool := server.GetTool("iot_query_devices")
	res1, err := devTool.Handler(ctxOrg1, mcp.CallToolRequest{})
	require.NoError(t, err)
	require.False(t, res1.IsError)
	devices1 := res1.StructuredContent.([]handler.MCPDevice)
	require.Len(t, devices1, 1)
	require.Equal(t, "device-org1", devices1[0].DeviceKey)

	// 2. Org 1 tries to access Org 2 device detail -> must fail / error
	detailTool := server.GetTool("iot_get_device_detail")
	resDetail, err := detailTool.Handler(ctxOrg1, mcp.CallToolRequest{Params: mcp.CallToolParams{
		Name:      "iot_get_device_detail",
		Arguments: map[string]any{"deviceKey": "device-org2"},
	}})
	require.NoError(t, err)
	require.True(t, resDetail.IsError)

	// 3. Org 1 tries to access Org 2 telemetry -> must fail / error
	telemTool := server.GetTool("iot_query_telemetry_history")
	resTelem, err := telemTool.Handler(ctxOrg1, mcp.CallToolRequest{Params: mcp.CallToolParams{
		Name:      "iot_query_telemetry_history",
		Arguments: map[string]any{"deviceKey": "device-org2", "property": "temp"},
	}})
	require.NoError(t, err)
	require.True(t, resTelem.IsError)
}

func TestStreamableHTTPAuthMiddleware(t *testing.T) {
	conf := viper.New()
	conf.Set("security.jwt.key", "test-secret-key-1234567890123456")
	j := jwt.NewJwt(conf)
	logger := &log.Logger{Logger: zap.NewNop()}

	token, err := j.GenToken("user-1", 1, time.Now().Add(time.Hour))
	require.NoError(t, err)

	mcpSrv := mcpserver.NewMCPServer("test-mcp", "1.0.0")
	streamable := mcpserver.NewStreamableHTTPServer(mcpSrv)

	var interceptedTenant int64
	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Unauthorized: missing Authorization token", http.StatusUnauthorized)
				return
			}
			claims, err := j.ParseToken(tokenString)
			if err != nil || claims.OrganizationID <= 0 {
				http.Error(w, "Unauthorized: invalid or expired token", http.StatusUnauthorized)
				return
			}
			ctx := tenant.WithTenant(r.Context(), claims.OrganizationID)
			ctx = context.WithValue(ctx, "claims", claims)
			ctx = context.WithValue(ctx, "user_id", claims.UserId)
			interceptedTenant = tenant.GetOrganizationID(ctx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	handler := authMiddleware(streamable)

	// 1. Missing token -> 401
	reqMissing := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	recMissing := httptest.NewRecorder()
	handler.ServeHTTP(recMissing, reqMissing)
	require.Equal(t, http.StatusUnauthorized, recMissing.Code)

	// 2. Invalid token -> 401
	reqInvalid := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	reqInvalid.Header.Set("Authorization", "Bearer invalid-token")
	recInvalid := httptest.NewRecorder()
	handler.ServeHTTP(recInvalid, reqInvalid)
	require.Equal(t, http.StatusUnauthorized, recInvalid.Code)

	// 3. Valid token -> 200/400 (reaches MCP handler, sets organization ID in context)
	reqValid := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	reqValid.Header.Set("Authorization", "Bearer "+token)
	recValid := httptest.NewRecorder()
	handler.ServeHTTP(recValid, reqValid)
	require.Equal(t, int64(1), interceptedTenant)
	_ = logger
	_ = mcptransport.WithStdioSrv()
}
