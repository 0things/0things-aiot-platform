package repository

import (
	"context"

	"data-engine/internal/model"
)

// DeviceEventRepository defines the storage interface for persisting device events.
type DeviceEventRepository interface {
	// Create persists a new device event record.
	Create(ctx context.Context, event *model.DeviceEvent) error
}

type deviceEventRepository struct {
	*Repository
}

// NewDeviceEventRepository initializes a new DeviceEventRepository instance with the base repository.
func NewDeviceEventRepository(repo *Repository) DeviceEventRepository {
	return &deviceEventRepository{Repository: repo}
}

// Create persists a device event record into PostgreSQL via GORM DAL.
func (r *deviceEventRepository) Create(ctx context.Context, event *model.DeviceEvent) error {
	if r.db == nil || event == nil {
		return nil
	}
	return useQuery(r.db).DeviceEvent.WithContext(ctx).Create(event)
}
