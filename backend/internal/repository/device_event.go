package repository

import (
	"context"
	"strings"
	"time"

	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"
	"gorm.io/gen/field"
)

type DeviceEventRepository struct{ db *IoTDB }

func NewDeviceEventRepository(db *IoTDB) *DeviceEventRepository {
	return &DeviceEventRepository{db: db}
}
func (r *DeviceEventRepository) Create(ctx context.Context, event *model.DeviceEvent) error {
	return useIoTQuery(r.db).DeviceEvent.WithContext(ctx).Create(event)
}

func (r *DeviceEventRepository) List(ctx context.Context, page, size int, keyword, deviceKey, eventType string, startAt, endAt *time.Time) ([]model.DeviceEvent, int64, error) {
	q := useIoTQuery(r.db)
	base := q.DeviceEvent.WithContext(ctx).Join(q.Device, q.Device.ID.EqCol(q.DeviceEvent.DeviceID)).Where(q.Device.OrganizationID.Eq(tenant.GetOrganizationID(ctx)))
	if keyword != "" {
		term := "%" + strings.ToLower(keyword) + "%"
		base = base.Where(field.Or(q.Device.DeviceKey.Lower().Like(term), q.Device.Name.Lower().Like(term), q.DeviceEvent.EventType.Lower().Like(term)))
	}
	if deviceKey != "" {
		base = base.Where(q.Device.DeviceKey.Eq(deviceKey))
	}
	if eventType != "" {
		base = base.Where(q.DeviceEvent.EventType.Lower().Like("%" + strings.ToLower(eventType) + "%"))
	}
	if startAt != nil {
		base = base.Where(q.DeviceEvent.EventAt.Gte(*startAt))
	}
	if endAt != nil {
		base = base.Where(q.DeviceEvent.EventAt.Lte(*endAt))
	}
	total, err := base.Count()
	if err != nil {
		return nil, 0, err
	}
	var events []model.DeviceEvent
	err = base.Select(q.DeviceEvent.ALL, q.Device.DeviceKey.As("device_key"), q.Device.Name.As("device_name")).Order(q.DeviceEvent.EventAt.Desc()).Offset((page - 1) * size).Limit(size).Scan(&events)
	return events, total, err
}
