package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"aiot-backend/internal/model"
	"gorm.io/gorm"
)

type DeviceShadowRepository struct {
	db *IoTDB
}

func NewDeviceShadowRepository(db *IoTDB) *DeviceShadowRepository {
	return &DeviceShadowRepository{db: db}
}

func (r *DeviceShadowRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *DeviceShadowRepository) GetShadow(ctx context.Context, deviceID int64) (*model.DeviceShadow, error) {
	var shadow model.DeviceShadow
	if err := r.DB(ctx).Where("device_id = ?", deviceID).First(&shadow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 影子是设备的按需资源，首次读取时初始化空文档，避免前端收到无意义的 404。
			shadow = model.DeviceShadow{DeviceID: deviceID, Desired: "{}", Reported: "{}", Metadata: "{}", Version: 0}
			if createErr := r.DB(ctx).Create(&shadow).Error; createErr != nil {
				return nil, createErr
			}
			return &shadow, nil
		}
		return nil, err
	}
	return &shadow, nil
}

// MutateShadow uses a version predicate in SQL so concurrent writes cannot
// silently overwrite each other. It also writes history in the same transaction.
func (r *DeviceShadowRepository) MutateShadow(ctx context.Context, deviceID, expected int64, source string, desired, reported *map[string]any, clearDesired bool) (*model.DeviceShadow, error) {
	var out model.DeviceShadow
	err := r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		var current model.DeviceShadow
		err := tx.Where("device_id = ?", deviceID).First(&current).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if expected != 0 {
				return ErrVersionConflict
			}
			current = model.DeviceShadow{DeviceID: deviceID, Desired: "{}", Reported: "{}", Metadata: "{}", Version: 0}
			if err := tx.Create(&current).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		if current.Version != expected {
			return ErrVersionConflict
		}
		if desired != nil {
			b, _ := json.Marshal(*desired)
			current.Desired = string(b)
		}
		if reported != nil {
			b, _ := json.Marshal(*reported)
			current.Reported = string(b)
		}
		if clearDesired {
			current.Desired = "{}"
		}
		current.Version++
		if err := tx.Model(&model.DeviceShadow{}).Where("id = ? AND version = ?", current.ID, expected).Updates(map[string]any{
			"desired": current.Desired, "reported": current.Reported, "version": current.Version, "updated_at": time.Now(),
		}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.DeviceShadowHistory{DeviceID: deviceID, Version: current.Version, Source: source, Desired: current.Desired, Reported: current.Reported, CreatedAt: time.Now()}).Error; err != nil {
			return err
		}
		out = current
		return nil
	})
	return &out, err
}

func (r *DeviceShadowRepository) ListShadowHistory(ctx context.Context, deviceID int64) ([]model.DeviceShadowHistory, error) {
	var history []model.DeviceShadowHistory
	err := r.DB(ctx).Where("device_id = ?", deviceID).Order("version DESC").Find(&history).Error
	return history, err
}
