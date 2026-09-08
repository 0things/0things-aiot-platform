package server

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"

	"transport-coap/pkg/log"
	"transport-coap/pkg/server"

	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	coapnet "github.com/plgd-dev/go-coap/v3/net"
	"github.com/plgd-dev/go-coap/v3/options"
	coapserver "github.com/plgd-dev/go-coap/v3/udp/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// CoAPServer wraps a UDP-based CoAP server for constrained device ingress.
type CoAPServer struct {
	addr   string
	server *coapserver.Server
	conn   *coapnet.UDPConn
	logger *log.Logger
	mu     sync.Mutex
}

var _ server.Server = (*CoAPServer)(nil)

// NewCoAPServer initializes CoAP UDP service.
func NewCoAPServer(config *viper.Viper, logger *log.Logger) *CoAPServer {
	addr := config.GetString("coap.addr")
	if addr == "" {
		addr = ":5683" // Standard CoAP port
	}
	return &CoAPServer{
		addr:   addr,
		logger: logger,
	}
}

// Start launches UDP listener and registers CoAP routes.
func (s *CoAPServer) Start(ctx context.Context) error {
	router := mux.NewRouter()

	router.Handle("/v1/device-ingress/{deviceKey}", mux.HandlerFunc(s.handleDeviceIngress))
	router.Handle("/api/v1/{deviceKey}/telemetry", mux.HandlerFunc(s.handleDeviceIngress))

	conn, err := coapnet.NewListenUDP("udp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen UDP on %s: %w", s.addr, err)
	}

	srv := coapserver.New(options.WithMux(router))

	s.mu.Lock()
	s.conn = conn
	s.server = srv
	s.mu.Unlock()

	s.logger.Info("CoAP Transport listening on UDP", zap.String("addr", s.addr))

	errChan := make(chan error, 1)
	go func() {
		if err := srv.Serve(conn); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return err
	}
}

// handleDeviceIngress processes a single CoAP datagram and responds with 2.04 Changed.
func (s *CoAPServer) handleDeviceIngress(w mux.ResponseWriter, r *mux.Message) {
	deviceKey := r.RouteParams.Vars["deviceKey"]
	if strings.TrimSpace(deviceKey) == "" {
		_ = w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("deviceKey is required")))
		return
	}

	payload, err := r.ReadBody()
	if err != nil || len(payload) == 0 {
		_ = w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("invalid body")))
		return
	}

	s.logger.Info("received CoAP ingress message",
		zap.String("device_key", deviceKey),
		zap.Int("payload_bytes", len(payload)),
	)

	// Reply CoAP 2.04 Changed so device can re-enter PSM power saving mode
	_ = w.SetResponse(codes.Changed, message.TextPlain, strings.NewReader("accepted"))
}

// Stop gracefully stops the CoAP server and releases the UDP socket.
func (s *CoAPServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.server != nil {
		s.server.Stop()
	}
	if s.conn != nil {
		_ = s.conn.Close()
	}
	s.logger.Info("CoAP server stopped")
	return nil
}
