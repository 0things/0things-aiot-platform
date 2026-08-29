package httptransport

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/transport"
)

// Adapter 接收设备 HTTP 上报，并转换为统一 DeviceMessage；业务认证和持久化由上层处理。
type Adapter struct {
	server *http.Server
}

func NewAdapter(addr string) *Adapter                { return &Adapter{server: &http.Server{Addr: addr}} }
func (a *Adapter) Name() string                      { return "http-device" }
func (a *Adapter) Transport() enum.TransportProtocol { return enum.TransportHTTP }

func (a *Adapter) Start(ctx context.Context, onMessage func(context.Context, transport.DeviceMessage) error) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/device-ingress/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		deviceKey := strings.TrimPrefix(r.URL.Path, "/v1/device-ingress/")
		if deviceKey == "" || strings.Contains(deviceKey, "/") {
			http.Error(w, "device key is required", http.StatusBadRequest)
			return
		}
		payload, err := io.ReadAll(io.LimitReader(r.Body, 4<<20))
		if err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		message := transport.DeviceMessage{DeviceKey: deviceKey, MessageType: r.Header.Get("X-Device-Message-Type"), Payload: payload, Timestamp: time.Now().UTC(), Headers: map[string]string{"content-type": r.Header.Get("Content-Type")}}
		if message.MessageType == "" {
			message.MessageType = "telemetry"
		}
		if onMessage != nil {
			if err := onMessage(r.Context(), message); err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
		}
		w.WriteHeader(http.StatusAccepted)
	})
	a.server.Handler = mux
	go func() { <-ctx.Done(); _ = a.Stop(context.Background()) }()
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *Adapter) Send(context.Context, transport.Command) error {
	return fmt.Errorf("http device downlink is not supported")
}
func (a *Adapter) Stop(ctx context.Context) error { return a.server.Shutdown(ctx) }
