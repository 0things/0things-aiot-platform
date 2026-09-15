package repository

import (
	"context"
	"strings"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// DeviceEventRepository provides database query operations for device events.
type DeviceEventRepository struct{ db *gorm.DB }

// NewDeviceEventRepository creates a new DeviceEventRepository instance.
func NewDeviceEventRepository(db *gorm.DB) *DeviceEventRepository {
	return &DeviceEventRepository{db: db}
}

// Create inserts a single device event record into the database.
func (r *DeviceEventRepository) Create(ctx context.Context, event *model.DeviceEvent) error {
	return useQuery(r.db).DeviceEvent.WithContext(ctx).Create(event)
}

// List queries device events with multi-tenant filtering, device joins, and pagination.
func (r *DeviceEventRepository) List(ctx context.Context, query dto.ListDeviceEventsQuery) ([]dto.DeviceEventListItem, int64, error) {
	q := useQuery(r.db)
	base := q.DeviceEvent.WithContext(ctx).Join(q.Device, q.Device.DeviceKey.EqCol(q.DeviceEvent.DeviceKey)).Where(q.Device.OrganizationID.Eq(tenant.GetOrganizationID(ctx)))
	if query.Keyword != "" {
		term := "%" + strings.ToLower(query.Keyword) + "%"
		base = base.Where(field.Or(q.Device.DeviceKey.Lower().Like(term), q.Device.Name.Lower().Like(term), q.DeviceEvent.EventType.Lower().Like(term)))
	}
	if query.DeviceKey != "" {
		base = base.Where(q.Device.DeviceKey.Eq(query.DeviceKey))
	}
	if query.EventType != "" {
		base = base.Where(q.DeviceEvent.EventType.Lower().Like("%" + strings.ToLower(query.EventType) + "%"))
	}
	if query.StartAt != nil {
		base = base.Where(q.DeviceEvent.EventAt.Gte(*query.StartAt))
	}
	if query.EndAt != nil {
		base = base.Where(q.DeviceEvent.EventAt.Lte(*query.EndAt))
	}
	total, err := base.Count()
	if err != nil {
		return nil, 0, err
	}
	var events []dto.DeviceEventListItem
	err = base.Select(q.DeviceEvent.ALL, q.Device.DeviceKey.As("device_key"), q.Device.Name.As("device_name")).Order(q.DeviceEvent.ID.Desc()).Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).Scan(&events)
	return events, total, err
}
