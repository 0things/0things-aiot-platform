package storage

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// DeviceShadow 描述设备的实时状态快照（设备影子）。
type DeviceShadow struct {
	DeviceKey  string                 `json:"device_key"`
	Attributes map[string]interface{} `json:"attributes"`
	LastSeen   time.Time              `json:"last_seen"`
}

// ShadowStore 设备影子存储接口。
type ShadowStore interface {
	UpdateShadow(ctx context.Context, deviceKey string, properties map[string]interface{}, ts time.Time) error
	GetShadow(ctx context.Context, deviceKey string) (*DeviceShadow, error)
}

// MemoryShadowStore 提供纯内存并发安全的设备影子实现（开箱即用无强依赖）。
type MemoryShadowStore struct {
	mu     sync.RWMutex
	shadow map[string]*DeviceShadow
	logger *zap.Logger
}

// NewShadowStore 初始化设备影子存储器。
func NewShadowStore(config *viper.Viper, logger *zap.Logger) *MemoryShadowStore {
	return &MemoryShadowStore{
		shadow: make(map[string]*DeviceShadow),
		logger: logger,
	}
}

// UpdateShadow 更新指定设备的影子属性快照。
func (s *MemoryShadowStore) UpdateShadow(ctx context.Context, deviceKey string, properties map[string]interface{}, ts time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.shadow[deviceKey]
	if !exists {
		current = &DeviceShadow{
			DeviceKey:  deviceKey,
			Attributes: make(map[string]interface{}),
		}
		s.shadow[deviceKey] = current
	}

	for k, v := range properties {
		current.Attributes[k] = v
	}
	current.LastSeen = ts

	s.logger.Debug("device shadow updated",
		zap.String("device_key", deviceKey),
		zap.Int("attr_count", len(current.Attributes)),
	)
	return nil
}

// GetShadow 获取单个设备的影子快照数据。
func (s *MemoryShadowStore) GetShadow(ctx context.Context, deviceKey string) (*DeviceShadow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if shadow, ok := s.shadow[deviceKey]; ok {
		// 返回只读拷贝
		attrsCopy := make(map[string]interface{}, len(shadow.Attributes))
		for k, v := range shadow.Attributes {
			attrsCopy[k] = v
		}
		return &DeviceShadow{
			DeviceKey:  shadow.DeviceKey,
			Attributes: attrsCopy,
			LastSeen:   shadow.LastSeen,
		}, nil
	}
	return nil, nil
}

// ToJSON 辅助序列化
func (d *DeviceShadow) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}
