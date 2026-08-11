package repository

import (
	"context"
	"errors"
	"time"

	"0things-backend/internal/model"
	"gorm.io/gorm"
)

type AlertRepository struct {
	db *IoTDB
}

func NewAlertRepository(db *IoTDB) *AlertRepository        { return &AlertRepository{db: db} }
func (r *AlertRepository) DB(ctx context.Context) *gorm.DB { return r.db.WithContext(ctx) }
func (r *AlertRepository) Find(ctx context.Context, id int64) (*model.Alert, error) {
	var alert model.Alert
	if err := r.DB(ctx).First(&alert, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &alert, nil
}

func (r *AlertRepository) List(ctx context.Context, page, size int, status, severity, deviceKey string) ([]model.Alert, int64, error) {
	query := r.DB(ctx).Model(&model.Alert{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if deviceKey != "" {
		query = query.Where("device_key = ?", deviceKey)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var alerts []model.Alert
	if err := query.Order("last_raised_at DESC").Offset((page - 1) * size).Limit(size).Find(&alerts).Error; err != nil {
		return nil, 0, err
	}
	return alerts, total, nil
}

func (r *AlertRepository) UpdateStatus(ctx context.Context, alert *model.Alert, status string, at time.Time) error {
	fields := map[string]any{"status": status}
	if status == "acknowledged" {
		fields["ack_at"] = at
	} else {
		fields["resolved_at"] = at
	}
	return r.DB(ctx).Model(alert).Updates(fields).Error
}
