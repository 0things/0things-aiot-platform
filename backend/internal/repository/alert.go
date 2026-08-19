package repository

import (
	"context"
	"errors"
	"time"

	"aiot-backend/internal/model"
	"gorm.io/gorm"
)

type AlertRepository struct {
	db *IoTDB
}

func NewAlertRepository(db *IoTDB) *AlertRepository { return &AlertRepository{db: db} }
func (r *AlertRepository) Find(ctx context.Context, id int64) (*model.Alert, error) {
	q := useIoTQuery(r.db)
	alert, err := q.Alert.WithContext(ctx).Where(q.Alert.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return alert, nil
}

func (r *AlertRepository) List(ctx context.Context, page, size int, status, severity, deviceKey string) ([]model.Alert, int64, error) {
	q := useIoTQuery(r.db)
	alerts := q.Alert.WithContext(ctx)
	if status != "" {
		alerts = alerts.Where(q.Alert.Status.Eq(status))
	}
	if severity != "" {
		alerts = alerts.Where(q.Alert.Severity.Eq(severity))
	}
	if deviceKey != "" {
		alerts = alerts.Where(q.Alert.DeviceKey.Eq(deviceKey))
	}
	items, total, err := alerts.Order(q.Alert.LastRaisedAt.Desc()).FindByPage((page-1)*size, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]model.Alert, len(items))
	for i := range items {
		result[i] = *items[i]
	}
	return result, total, nil
}

func (r *AlertRepository) UpdateStatus(ctx context.Context, alert *model.Alert, status string, at time.Time) error {
	q := useIoTQuery(r.db)
	if status == "acknowledged" {
		_, err := q.Alert.WithContext(ctx).Where(q.Alert.ID.Eq(alert.ID)).UpdateSimple(q.Alert.Status.Value(status), q.Alert.AckAt.Value(at))
		return err
	}
	_, err := q.Alert.WithContext(ctx).Where(q.Alert.ID.Eq(alert.ID)).UpdateSimple(q.Alert.Status.Value(status), q.Alert.ResolvedAt.Value(at))
	return err
}
