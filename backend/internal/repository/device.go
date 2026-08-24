package repository

import (
	"context"
	"errors"
	"fmt"

	"aiot-backend/internal/dal/query"
	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"
	"github.com/redis/go-redis/v9"
	"gorm.io/gen/field"
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
	q := useQuery(r.db)
	item, err := q.Device.WithContext(ctx).Where(q.Device.ID.Eq(id), q.Device.OrganizationID.Eq(tenant.GetOrganizationID(ctx))).Preload(q.Device.Product, q.Device.State).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}

func (r *DeviceRepository) FindByKey(ctx context.Context, key string) (*model.Device, error) {
	q := useQuery(r.db)
	item, err := q.Device.WithContext(ctx).Where(q.Device.DeviceKey.Eq(key), q.Device.OrganizationID.Eq(tenant.GetOrganizationID(ctx))).Preload(q.Device.Product, q.Device.State).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}

// FindByKeys resolves multiple devices by their device keys in a single query.
func (r *DeviceRepository) FindByKeys(ctx context.Context, keys []string) ([]*model.Device, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	q := useQuery(r.db)
	devices, err := q.Device.WithContext(ctx).
		Where(q.Device.DeviceKey.In(keys...), q.Device.OrganizationID.Eq(tenant.GetOrganizationID(ctx))).
		Find()
	if err != nil {
		return nil, err
	}
	return devices, nil
}

// FindByKeyForEvent resolves a device from a broker event, which has no user
// request context. Device keys are globally unique in the device protocol.
func (r *DeviceRepository) FindByKeyForEvent(ctx context.Context, key string) (*model.Device, error) {
	q := useQuery(r.db)
	item, err := q.Device.WithContext(ctx).Where(q.Device.DeviceKey.Eq(key)).Preload(q.Device.Product, q.Device.State).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}

func (r *DeviceRepository) Create(ctx context.Context, device *model.Device) error {
	return useQuery(r.db).Transaction(func(tx *query.Query) error {
		if err := tx.Device.WithContext(ctx).Create(device); err != nil {
			return err
		}
		return tx.DeviceState.WithContext(ctx).Create(&model.DeviceState{DeviceKey: device.DeviceKey, State: "inactive"})
	})
}

func (r *DeviceRepository) List(ctx context.Context, page, size int, productID int64, states []string, enabled *bool, search string) ([]model.Device, int64, error) {
	q := useQuery(r.db)
	devices := q.Device.WithContext(ctx).Join(q.DeviceState, q.DeviceState.DeviceKey.EqCol(q.Device.DeviceKey)).Where(q.Device.OrganizationID.Eq(tenant.GetOrganizationID(ctx)))
	if productID > 0 {
		devices = devices.Where(q.Device.ProductID.Eq(productID))
	}
	if len(states) > 0 {
		devices = devices.Where(q.DeviceState.State.In(states...))
	}
	if enabled != nil {
		devices = devices.Where(q.Device.Enabled.Is(*enabled))
	}
	if search != "" {
		devices = devices.Where(field.Or(q.Device.DeviceKey.Like("%"+search+"%"), q.Device.Name.Like("%"+search+"%")))
	}
	items, total, err := devices.Preload(q.Device.Product, q.Device.State).Order(q.Device.CreatedAt.Desc()).FindByPage((page-1)*size, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]model.Device, len(items))
	for i := range items {
		result[i] = *items[i]
	}
	return result, total, nil
}

func (r *DeviceRepository) Save(ctx context.Context, device *model.Device) error {
	return useQuery(r.db).Transaction(func(tx *query.Query) error {
		if err := tx.Device.WithContext(ctx).Save(device); err != nil {
			return err
		}
		return tx.DeviceState.WithContext(ctx).Save(&device.State)
	})
}

func (r *DeviceRepository) SaveEnabled(ctx context.Context, device *model.Device) error {
	return useQuery(r.db).Device.WithContext(ctx).Save(device)
}

func (r *DeviceRepository) Statistics(ctx context.Context) (DeviceStatistics, error) {
	var total, online, offline, inactive int64
	var err error
	q := useQuery(r.db)
	if total, err = q.Device.WithContext(ctx).Where(q.Device.OrganizationID.Eq(tenant.GetOrganizationID(ctx))).Count(); err != nil {
		return DeviceStatistics{}, err
	}
	for state, target := range map[string]*int64{"online": &online, "offline": &offline, "inactive": &inactive} {
		count, err := q.Device.WithContext(ctx).Join(q.DeviceState, q.DeviceState.DeviceKey.EqCol(q.Device.DeviceKey)).Where(q.Device.OrganizationID.Eq(tenant.GetOrganizationID(ctx)), q.DeviceState.State.Eq(state)).Count()
		*target = count
		if err != nil {
			return DeviceStatistics{}, err
		}
	}
	return DeviceStatistics{
		TotalDevices: total, ActivatedDevices: online + offline, OnlineDevices: online,
		OfflineDevices: offline, InactiveDevices: inactive,
	}, nil
}

func (r *DeviceRepository) Delete(ctx context.Context, device *model.Device) error {
	q := useQuery(r.db)
	_, err := q.Device.WithContext(ctx).Where(q.Device.OrganizationID.Eq(tenant.GetOrganizationID(ctx))).Delete(device)
	return err
}

func (r *DeviceRepository) Restore(ctx context.Context, id int64) error {
	q := useQuery(r.db)
	_, err := q.Device.WithContext(ctx).Unscoped().Where(q.Device.ID.Eq(id), q.Device.OrganizationID.Eq(tenant.GetOrganizationID(ctx))).UpdateSimple(q.Device.DeletedAt.Null())
	return err
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
