package repository

import (
	"context"
	"errors"
	"fmt"

	"0things-backend/internal/model"
	"0things-backend/internal/repository/scope"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrVersionConflict = errors.New("shadow version conflict")
)

type DeviceRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

type DeviceStatistics struct {
	TotalDevices     int64
	ActivatedDevices int64
	OnlineDevices    int64
	OfflineDevices   int64
	InactiveDevices  int64
}

func NewDeviceRepository(db *IoTDB, redis *IoTRedis) *DeviceRepository {
	return &DeviceRepository{db: db.DB, redis: redis.Client}
}

func (r *DeviceRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *DeviceRepository) Find(ctx context.Context, id int64) (*model.Device, error) {
	var item model.Device
	if err := r.DB(ctx).Scopes(scope.Tenant).Preload("Product").Preload("State").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *DeviceRepository) FindByKey(ctx context.Context, key string) (*model.Device, error) {
	var item model.Device
	if err := r.DB(ctx).Scopes(scope.Tenant).Preload("Product").Preload("State").Where("device_key = ?", key).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

// FindByKeyForEvent resolves a device from a broker event, which has no user
// request context. Device keys are globally unique in the device protocol.
func (r *DeviceRepository) FindByKeyForEvent(ctx context.Context, key string) (*model.Device, error) {
	var item model.Device
	if err := r.DB(ctx).Preload("Product").Preload("State").Where("device_key = ?", key).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *DeviceRepository) Create(ctx context.Context, device *model.Device) error {
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(device).Error; err != nil {
			return err
		}
		return tx.Create(&model.DeviceState{DeviceID: device.ID, State: "inactive"}).Error
	})
}

func (r *DeviceRepository) List(ctx context.Context, page, size int, productID int64, states []string, enabled *bool, search string) ([]model.Device, int64, error) {
	query := r.DB(ctx).Model(&model.Device{}).Scopes(scope.Tenant).Joins("JOIN device_states ON device_states.device_id = devices.id")
	if productID > 0 {
		query = query.Where("devices.product_id = ?", productID)
	}
	if len(states) > 0 {
		query = query.Where("device_states.state IN ?", states)
	}
	if enabled != nil {
		query = query.Where("devices.enabled = ?", *enabled)
	}
	if search != "" {
		query = query.Where("devices.device_key LIKE ? OR devices.name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var devices []model.Device
	if err := query.Preload("Product").Preload("State").Order("devices.created_at DESC").Offset((page - 1) * size).Limit(size).Find(&devices).Error; err != nil {
		return nil, 0, err
	}
	return devices, total, nil
}

func (r *DeviceRepository) Save(ctx context.Context, device *model.Device) error {
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(device).Error; err != nil {
			return err
		}
		return tx.Save(&device.State).Error
	})
}

func (r *DeviceRepository) SaveEnabled(ctx context.Context, device *model.Device) error {
	return r.DB(ctx).Save(device).Error
}

func (r *DeviceRepository) Statistics(ctx context.Context) (DeviceStatistics, error) {
	var total, online, offline, inactive int64
	if err := r.DB(ctx).Scopes(scope.Tenant).Model(&model.Device{}).Count(&total).Error; err != nil {
		return DeviceStatistics{}, err
	}
	for state, target := range map[string]*int64{"online": &online, "offline": &offline, "inactive": &inactive} {
		if err := r.DB(ctx).Scopes(scope.Tenant).Model(&model.Device{}).Joins("JOIN device_states ON device_states.device_id = devices.id").Where("device_states.state = ?", state).Count(target).Error; err != nil {
			return DeviceStatistics{}, err
		}
	}
	return DeviceStatistics{
		TotalDevices: total, ActivatedDevices: online + offline, OnlineDevices: online,
		OfflineDevices: offline, InactiveDevices: inactive,
	}, nil
}

func (r *DeviceRepository) Delete(ctx context.Context, device *model.Device) error {
	return r.DB(ctx).Delete(device).Error
}

func (r *DeviceRepository) Restore(ctx context.Context, id int64) error {
	return r.DB(ctx).Unscoped().Scopes(scope.Tenant).Model(&model.Device{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *DeviceRepository) Telemetry(ctx context.Context, key string) (string, error) {
	if r.redis == nil {
		return "", fmt.Errorf("redis is not configured")
	}
	v, err := r.redis.Get(ctx, "telemetry:device:"+key+":latest").Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return v, err
}
