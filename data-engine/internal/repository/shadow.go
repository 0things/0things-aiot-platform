package repository

import (
	"context"
	"sync"
	"time"

	"data-engine/internal/model"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ShadowRepository defines the storage interface for device shadow snapshots.
type ShadowRepository interface {
	UpdateShadow(ctx context.Context, deviceKey string, properties map[string]interface{}, ts time.Time) error
	GetShadow(ctx context.Context, deviceKey string) (*model.DeviceShadow, error)
}

// memoryShadowRepository provides an in-memory concurrent-safe device shadow implementation.
type memoryShadowRepository struct {
	mu     sync.RWMutex
	shadow map[string]*model.DeviceShadow
	logger *zap.Logger
}

// NewShadowRepository initializes the device shadow repository.
func NewShadowRepository(config *viper.Viper, logger *zap.Logger) ShadowRepository {
	return &memoryShadowRepository{
		shadow: make(map[string]*model.DeviceShadow),
		logger: logger,
	}
}

// UpdateShadow updates property snapshot for a device shadow.
func (s *memoryShadowRepository) UpdateShadow(ctx context.Context, deviceKey string, properties map[string]interface{}, ts time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.shadow[deviceKey]
	if !exists {
		current = &model.DeviceShadow{
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

// GetShadow fetches the current device shadow snapshot.
func (s *memoryShadowRepository) GetShadow(ctx context.Context, deviceKey string) (*model.DeviceShadow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if shadow, ok := s.shadow[deviceKey]; ok {
		attrsCopy := make(map[string]interface{}, len(shadow.Attributes))
		for k, v := range shadow.Attributes {
			attrsCopy[k] = v
		}
		return &model.DeviceShadow{
			DeviceKey:  shadow.DeviceKey,
			Attributes: attrsCopy,
			LastSeen:   shadow.LastSeen,
		}, nil
	}
	return nil, nil
}
