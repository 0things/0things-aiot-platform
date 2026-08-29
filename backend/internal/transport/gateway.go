package transport

import (
	"context"
	"sync"
)

// Gateway 负责统一启动和停止多个协议适配器；适配器之间互不持有业务服务。
type Gateway struct {
	adapters  []Adapter
	onMessage func(context.Context, DeviceMessage) error
	mu        sync.Mutex
	started   bool
}

func NewGateway(adapters []Adapter, onMessage func(context.Context, DeviceMessage) error) *Gateway {
	return &Gateway{adapters: adapters, onMessage: onMessage}
}

func (g *Gateway) Start(ctx context.Context) error {
	g.mu.Lock()
	if g.started {
		g.mu.Unlock()
		return nil
	}
	g.started = true
	g.mu.Unlock()
	for _, adapter := range g.adapters {
		go func(a Adapter) {
			if err := a.Start(ctx, g.onMessage); err != nil && ctx.Err() == nil {
				// 适配器错误由具体进程的日志和监控接管，不能阻塞其他协议启动。
			}
		}(adapter)
	}
	return nil
}

func (g *Gateway) Stop(ctx context.Context) error {
	g.mu.Lock()
	if !g.started {
		g.mu.Unlock()
		return nil
	}
	g.started = false
	g.mu.Unlock()
	for i := len(g.adapters) - 1; i >= 0; i-- {
		if err := g.adapters[i].Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}
