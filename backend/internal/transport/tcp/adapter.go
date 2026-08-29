package tcptransport

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"aiot-backend/internal/enum"
	"aiot-backend/internal/transport"
)

// Adapter 提供最小 TCP 设备接入：每行格式为 deviceKey|payload，业务协议由 Codec 继续解析。
type Adapter struct {
	addr     string
	listener net.Listener
	mu       sync.Mutex
}

func NewAdapter(addr string) *Adapter                { return &Adapter{addr: addr} }
func (a *Adapter) Name() string                      { return "tcp-device" }
func (a *Adapter) Transport() enum.TransportProtocol { return enum.TransportTCP }

func (a *Adapter) Start(ctx context.Context, onMessage func(context.Context, transport.DeviceMessage) error) error {
	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.listener = listener
	a.mu.Unlock()
	go func() { <-ctx.Done(); _ = a.Stop(context.Background()) }()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		go a.handleConnection(ctx, conn, onMessage)
	}
}

func (a *Adapter) handleConnection(ctx context.Context, conn net.Conn, onMessage func(context.Context, transport.DeviceMessage) error) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "|", 2)
		if len(parts) != 2 || parts[0] == "" {
			continue
		}
		if onMessage != nil {
			_ = onMessage(ctx, transport.DeviceMessage{DeviceKey: parts[0], MessageType: "telemetry", Payload: []byte(parts[1]), Timestamp: time.Now().UTC(), Headers: map[string]string{"transport": "tcp"}})
		}
	}
}

func (a *Adapter) Send(ctx context.Context, command transport.Command) error {
	endpoint := command.Headers["endpoint"]
	if endpoint == "" {
		return fmt.Errorf("tcp endpoint is required")
	}
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", endpoint)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = fmt.Fprintf(conn, "%s|%s\n", command.DeviceKey, command.Payload)
	return err
}

func (a *Adapter) Stop(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.listener == nil {
		return nil
	}
	return a.listener.Close()
}
