package coap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"coap-transport/internal/kafka"
	"coap-transport/internal/model"

	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	coapnet "github.com/plgd-dev/go-coap/v3/net"
	"github.com/plgd-dev/go-coap/v3/options"
	coapserver "github.com/plgd-dev/go-coap/v3/udp/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Service 封装基于 UDP 的 CoAP 服务端，专职接收 NB-IoT 与受限低功耗设备上报。
type Service struct {
	addr     string
	producer *kafka.Producer
	server   *coapserver.Server
	conn     *coapnet.UDPConn
	logger   *zap.Logger
	mu       sync.Mutex
}

// NewService 初始化 CoAP UDP 监听服务。
func NewService(config *viper.Viper, logger *zap.Logger, producer *kafka.Producer) *Service {
	addr := config.GetString("coap.addr")
	if addr == "" {
		addr = ":5683" // CoAP 标准默认端口
	}
	return &Service{
		addr:     addr,
		producer: producer,
		logger:   logger,
	}
}

// Start 启动 UDP 监听并挂载 CoAP 资源路由。
func (s *Service) Start(ctx context.Context) error {
	router := mux.NewRouter()

	// 注册 CoAP 设备上报路径（支持老网关路径与标准 v1 路径）
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

// handleDeviceIngress 处理单个 CoAP 报文，提取 deviceKey 与 Body 后封装入 Kafka，并向设备回复 CoAP 2.04 Changed 响应。
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

	msg := model.DeviceMessage{
		DeviceKey:   deviceKey,
		Transport:   "coap",
		MessageType: "telemetry",
		Payload:     json.RawMessage(payload),
		Timestamp:   time.Now().UTC(),
		Headers: func() map[string]string {
			h := make(map[string]string)
			if addrVal := r.Context().Value("remote-addr"); addrVal != nil {
				if s, ok := addrVal.(fmt.Stringer); ok {
					h["remote-addr"] = s.String()
				}
			}
			return h
		}(),
	}

	// 投递 Kafka
	if err := s.producer.SendDeviceMessage(r.Context(), msg); err != nil {
		s.logger.Error("failed to produce coap message to kafka", zap.String("device_key", deviceKey), zap.Error(err))
		_ = w.SetResponse(codes.BadGateway, message.TextPlain, strings.NewReader("queue error"))
		return
	}

	// 向低功耗终端回复标准 CoAP 2.04 Changed 状态码，使设备能够快速重新进入深度休眠（PSM）
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
