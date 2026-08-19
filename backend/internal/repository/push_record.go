package repository

import (
	"context"
	"errors"
	"time"

	"aiot-backend/internal/model"
	"gorm.io/gorm"
)

type PushRecordRepository struct {
	db *IoTDB
}

func NewPushRecordRepository(db *IoTDB) *PushRecordRepository {
	return &PushRecordRepository{db: db}
}

func (r *PushRecordRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *PushRecordRepository) CreatePushRecord(ctx context.Context, record *model.DevicePushRecord) error {
	return r.DB(ctx).Create(record).Error
}

func (r *PushRecordRepository) ListPushRecords(ctx context.Context, deviceID int64, page, size int, operationType, status string) ([]model.DevicePushRecord, int64, error) {
	query := r.DB(ctx).Model(&model.DevicePushRecord{}).Where("device_id = ?", deviceID)
	if operationType != "" {
		query = query.Where("operation_type = ?", operationType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var records []model.DevicePushRecord
	if err := query.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

func (r *PushRecordRepository) FindPushRecord(ctx context.Context, id int64) (*model.DevicePushRecord, error) {
	var record model.DevicePushRecord
	if err := r.DB(ctx).First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &record, nil
}

func (r *PushRecordRepository) DeletePushRecords(ctx context.Context, deviceID int64, before *time.Time) (int64, error) {
	query := r.DB(ctx).Where("device_id = ?", deviceID)
	if before != nil {
		query = query.Where("created_at < ?", *before)
	}
	result := query.Delete(&model.DevicePushRecord{})
	return result.RowsAffected, result.Error
}
