package server

import (
	v1 "telemetry-service/api/helloworld/v1"
	clientConnectStatev1 "telemetry-service/api/client_connect_state/v1"
	eventv1 "telemetry-service/api/event/v1"
	telemetryv1 "telemetry-service/api/telemetry/v1"
	"telemetry-service/internal/conf"
	"telemetry-service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, telemetry *service.TelemetryService, event *service.EventService, clientConnectState *service.ClientConnectStateService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	telemetryv1.RegisterTelemetryServiceHTTPServer(srv, telemetry)
	eventv1.RegisterEventServiceHTTPServer(srv, event)
	clientConnectStatev1.RegisterClientConnectStateServiceHTTPServer(srv, clientConnectState)
	return srv
}
