package repository

import (
	"context"
	"errors"

	"aiot-backend/internal/model"

	"gorm.io/gorm"
)

type DeviceTagRepository struct {
	db *gorm.DB
}

func NewDeviceTagRepository(db *gorm.DB) *DeviceTagRepository {
	return &DeviceTagRepository{db: db}
}

func (r *DeviceTagRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *DeviceTagRepository) ListTags(ctx context.Context, deviceID int64) ([]model.DeviceTag, error) {
	var tags []model.DeviceTag
	err := r.DB(ctx).Where("device_id = ?", deviceID).Order("key").Find(&tags).Error
	return tags, err
}

func (r *DeviceTagRepository) SetTags(ctx context.Context, deviceID int64, tags map[string]string, replace bool) error {
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if replace {
			if err := tx.Where("device_id = ?", deviceID).Delete(&model.DeviceTag{}).Error; err != nil {
				return err
			}
		}
		for key, value := range tags {
			var tag model.DeviceTag
			err := tx.Unscoped().Where("device_id = ? AND key = ?", deviceID, key).First(&tag).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := tx.Create(&model.DeviceTag{DeviceID: deviceID, Key: key, Value: value, Source: "manual"}).Error; err != nil {
					return err
				}
				continue
			}
			if err != nil {
				return err
			}
			if err := tx.Unscoped().Model(&tag).Updates(map[string]any{"value": value, "source": "manual", "deleted_at": nil}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *DeviceTagRepository) DeleteTags(ctx context.Context, deviceID int64, keys []string) error {
	return r.DB(ctx).Where("device_id = ? AND key IN ?", deviceID, keys).Delete(&model.DeviceTag{}).Error
}
