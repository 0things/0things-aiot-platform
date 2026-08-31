package protocol

import (
	"fmt"
	"sync"

	gb28181codec "0things/pkg/protocol/gb28181"
	jsoncodec "0things/pkg/protocol/json"
	modbuscodec "0things/pkg/protocol/modbus"
)

// Registry 统一管理应用层协议编解码器，支持即插即用扩展。
type Registry struct {
	mu     sync.RWMutex
	codecs map[string]Codec
}

func NewRegistry(codecs ...Codec) (*Registry, error) {
	r := &Registry{codecs: make(map[string]Codec)}
	for _, codec := range codecs {
		if codec == nil {
			continue
		}
		if _, exists := r.codecs[codec.Name()]; exists {
			return nil, fmt.Errorf("application codec %q already registered", codec.Name())
		}
		r.codecs[codec.Name()] = codec
	}
	return r, nil
}

// DefaultRegistry 返回预装 JSON / Modbus / GB28181 常用编解码器的全局注册表
func DefaultRegistry() *Registry {
	r, _ := NewRegistry(
		jsoncodec.New(),
		modbuscodec.New(),
		gb28181codec.New(),
	)
	return r
}

func (r *Registry) Register(codec Codec) error {
	if codec == nil {
		return fmt.Errorf("codec cannot be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.codecs[codec.Name()]; exists {
		return fmt.Errorf("codec %q already registered", codec.Name())
	}
	r.codecs[codec.Name()] = codec
	return nil
}

func (r *Registry) Get(name string) (Codec, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	codec, ok := r.codecs[name]
	return codec, ok
}
