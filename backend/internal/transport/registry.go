package transport

import (
	"fmt"
	"sync"

	"aiot-backend/internal/enum"
)

// Registry 统一管理协议适配器，设备和 OTA 服务无需编写协议分支。
type Registry struct {
	mu       sync.RWMutex
	adapters map[enum.TransportProtocol]Adapter
}

func NewRegistry(adapters ...Adapter) (*Registry, error) {
	r := &Registry{adapters: make(map[enum.TransportProtocol]Adapter)}
	for _, adapter := range adapters {
		if adapter == nil {
			continue
		}
		if _, exists := r.adapters[adapter.Transport()]; exists {
			return nil, fmt.Errorf("transport adapter %q already registered", adapter.Transport())
		}
		r.adapters[adapter.Transport()] = adapter
	}
	return r, nil
}

func (r *Registry) Get(protocol enum.TransportProtocol) (Adapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.adapters[protocol]
	return a, ok
}
