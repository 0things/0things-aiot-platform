package coaptransport

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/transport"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	coapnet "github.com/plgd-dev/go-coap/v3/net"
	"github.com/plgd-dev/go-coap/v3/options"
	coapudp "github.com/plgd-dev/go-coap/v3/udp"
	coapserver "github.com/plgd-dev/go-coap/v3/udp/server"
)

// Adapter 接收标准 CoAP POST，上行路径使用 /v1/device-ingress/{deviceKey}。
type Adapter struct {
	addr   string
	server *coapserver.Server
	conn   *coapnet.UDPConn
	mu     sync.Mutex
}

func NewAdapter(addr string) *Adapter                { return &Adapter{addr: addr} }
func (a *Adapter) Name() string                      { return "coap-device" }
func (a *Adapter) Transport() enum.TransportProtocol { return enum.TransportCoAP }

func (a *Adapter) Start(ctx context.Context, onMessage func(context.Context, transport.DeviceMessage) error) error {
	router := mux.NewRouter()
	router.Handle("/v1/device-ingress/{deviceKey}", mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		deviceKey := r.RouteParams.Vars["deviceKey"]
		payload, err := r.ReadBody()
		if err != nil || strings.TrimSpace(deviceKey) == "" {
			_ = w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("invalid device message")))
			return
		}
		if onMessage != nil {
			if err = onMessage(r.Context(), transport.DeviceMessage{DeviceKey: deviceKey, MessageType: "telemetry", Payload: payload, Timestamp: time.Now().UTC(), Headers: map[string]string{"transport": "coap"}}); err != nil {
				_ = w.SetResponse(codes.BadGateway, message.TextPlain, strings.NewReader(err.Error()))
				return
			}
		}
		_ = w.SetResponse(codes.Changed, message.TextPlain, strings.NewReader("accepted"))
	}))
	conn, err := coapnet.NewListenUDP("udp", a.addr)
	if err != nil {
		return err
	}
	server := coapserver.New(options.WithMux(router))
	a.mu.Lock()
	a.conn, a.server = conn, server
	a.mu.Unlock()
	go func() { <-ctx.Done(); _ = a.Stop(context.Background()) }()
	err = server.Serve(conn)
	if ctx.Err() != nil {
		return nil
	}
	return err
}

func (a *Adapter) Send(ctx context.Context, command transport.Command) error {
	endpoint := command.Headers["endpoint"]
	if endpoint == "" {
		return fmt.Errorf("coap endpoint is required")
	}
	endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "coap://"), "coaps://")
	conn, err := coapudp.Dial(endpoint)
	if err != nil {
		return err
	}
	defer conn.Close()
	path := command.Headers["path"]
	if path == "" {
		path = "/v1/device-upgrade/" + command.DeviceKey
	}
	resp, err := conn.Post(ctx, path, message.AppJSON, bytes.NewReader(command.Payload))
	if err != nil {
		return err
	}
	_, err = io.Copy(io.Discard, resp.Body())
	return err
}

func (a *Adapter) Stop(context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.server != nil {
		a.server.Stop()
		a.server = nil
	}
	if a.conn != nil {
		err := a.conn.Close()
		a.conn = nil
		return err
	}
	return nil
}
