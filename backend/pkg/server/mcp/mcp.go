package mcp

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	mcpserver "github.com/mark3labs/mcp-go/server"
)

type transport struct {
	start func(context.Context) error
	stop  func(context.Context) error
}

// Server adapts MCP transports to the application's Start/Stop lifecycle.
type Server struct {
	transports []transport
	cancel     context.CancelFunc
	mu         sync.Mutex
}

type Option func(*Server, *mcpserver.MCPServer)

type HTTPMiddleware func(http.Handler) http.Handler

func NewServer(server *mcpserver.MCPServer, opts ...Option) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s, server)
	}
	return s
}

func WithStdioSrv() Option {
	return func(s *Server, server *mcpserver.MCPServer) {
		stdio := mcpserver.NewStdioServer(server)
		s.transports = append(s.transports, transport{
			start: func(ctx context.Context) error {
				return stdio.Listen(ctx, os.Stdin, os.Stdout)
			},
			stop: func(context.Context) error { return nil },
		})
	}
}

func WithSSESrv(addr string) Option {
	return func(s *Server, server *mcpserver.MCPServer) {
		sse := mcpserver.NewSSEServer(server)
		s.transports = append(s.transports, transport{
			start: func(context.Context) error { return sse.Start(addr) },
			stop:  sse.Shutdown,
		})
	}
}

func WithStreamableHTTPSrv(addr string, middlewares ...HTTPMiddleware) Option {
	return func(s *Server, server *mcpserver.MCPServer) {
		streamable := mcpserver.NewStreamableHTTPServer(server)
		var h http.Handler = streamable
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		httpSrv := &http.Server{
			Addr:    addr,
			Handler: h,
		}
		s.transports = append(s.transports, transport{
			start: func(ctx context.Context) error {
				go func() {
					<-ctx.Done()
					shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					_ = httpSrv.Shutdown(shutdownCtx)
				}()
				if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					return err
				}
				return nil
			},
			stop: func(ctx context.Context) error {
				return httpSrv.Shutdown(ctx)
			},
		})
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.cancel != nil {
		s.mu.Unlock()
		return errors.New("MCP server is already running")
	}
	if len(s.transports) == 0 {
		s.mu.Unlock()
		return errors.New("MCP server has no configured transport")
	}

	runCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.mu.Unlock()

	errs := make(chan error, len(s.transports))
	for _, item := range s.transports {
		go func(item transport) {
			if err := item.start(runCtx); err != nil && !errors.Is(err, context.Canceled) {
				errs <- err
			}
		}(item)
	}

	select {
	case err := <-errs:
		return err
	case <-runCtx.Done():
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	cancel := s.cancel
	s.cancel = nil
	transports := append([]transport(nil), s.transports...)
	s.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	var errs []error
	for _, transport := range transports {
		if err := transport.stop(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
