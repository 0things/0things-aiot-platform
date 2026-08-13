package repository

import (
	"context"
	"strings"
	"time"

	"0things-backend/internal/model"
	"gorm.io/gorm"
)

type DeviceEventRepository struct{ db *IoTDB }

func NewDeviceEventRepository(db *IoTDB) *DeviceEventRepository { return &DeviceEventRepository{db: db} }
func (r *DeviceEventRepository) DB(ctx context.Context) *gorm.DB { return r.db.WithContext(ctx) }

func (r *DeviceEventRepository) Create(ctx context.Context, event *model.DeviceEvent) error {
	return r.DB(ctx).Create(event).Error
}

func (r *DeviceEventRepository) List(ctx context.Context, page, size int, keyword, deviceKey, eventType string, startAt, endAt *time.Time) ([]model.DeviceEvent, int64, error) {
	base := r.DB(ctx).Model(&model.DeviceEvent{}).
		Joins("LEFT JOIN devices d ON d.id = device_events.device_id")
	if keyword != "" {
		term := "%" + strings.ToLower(keyword) + "%"
		base = base.Where("LOWER(d.device_key) LIKE ? OR LOWER(d.name) LIKE ? OR LOWER(device_events.event_type) LIKE ?", term, term, term)
	}
	if deviceKey != "" { base = base.Where("d.device_key = ?", deviceKey) }
	if eventType != "" {
		base = base.Where("LOWER(device_events.event_type) LIKE ?", "%"+strings.ToLower(eventType)+"%")
	}
	if startAt != nil { base = base.Where("device_events.event_at >= ?", *startAt) }
	if endAt != nil { base = base.Where("device_events.event_at <= ?", *endAt) }
	var total int64
	if err := base.Count(&total).Error; err != nil { return nil, 0, err }
	var events []model.DeviceEvent
	err := base.Select("device_events.*, d.device_key AS device_key, d.name AS device_name").
		Order("device_events.event_at DESC").Offset((page-1)*size).Limit(size).Find(&events).Error
	return events, total, err
}
