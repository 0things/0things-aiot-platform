package coap

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	coapnet "github.com/plgd-dev/go-coap/v3/net"
	"github.com/plgd-dev/go-coap/v3/options"
	coapserver "github.com/plgd-dev/go-coap/v3/udp/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Service wraps a UDP-based CoAP server for constrained device ingress.
type Service struct {
	addr   string
	server *coapserver.Server
	conn   *coapnet.UDPConn
	logger *zap.Logger
	mu     sync.Mutex
}

// NewService initializes CoAP UDP service.
func NewService(config *viper.Viper, logger *zap.Logger) *Service {
	addr := config.GetString("coap.addr")
	if addr == "" {
		addr = ":5683" // Standard CoAP port
	}
	return &Service{
		addr:   addr,
		logger: logger,
	}
}

// Start launches UDP listener and registers CoAP routes.
func (s *Service) Start(ctx context.Context) error {
	router := mux.NewRouter()

	router.Handle("/v1/device-ingress/{deviceKey}", mux.HandlerFunc(s.handleDeviceIngress))
	router.Handle("/api/v1/{deviceKey}/telemetry", mux.HandlerFunc(s.handleDeviceIngress))

	conn, err := coapnet.NewListenUDP("udp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen UDP on %s: %w", s.addr, err)
	}

	server := coapserver.New(options.WithMux(router))

	s.mu.Lock()
	s.conn = conn
	s.server = server
	s.mu.Unlock()

	s.logger.Info("CoAP Transport listening on UDP", zap.String("addr", s.addr))

	errChan := make(chan error, 1)
	go func() {
		if err := server.Serve(conn); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.Stop()
		return nil
	case err := <-errChan:
		return err
	}
}

// handleDeviceIngress processes a single CoAP datagram and responds with 2.04 Changed.
func (s *Service) handleDeviceIngress(w mux.ResponseWriter, r *mux.Message) {
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

// Stop 安全停止 CoAP 服务端并释放 UDP Socket
func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.server != nil {
		s.server.Stop()
	}
	if s.conn != nil {
		_ = s.conn.Close()
	}
	s.logger.Info("CoAP server stopped")
}
