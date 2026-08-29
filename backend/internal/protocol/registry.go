package protocol

import (
	"fmt"
	"sync"
)

// Registry 统一管理应用协议编解码器，新增协议无需修改设备和 OTA 业务。
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

func (r *Registry) Get(name string) (Codec, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	codec, ok := r.codecs[name]
	return codec, ok
}
