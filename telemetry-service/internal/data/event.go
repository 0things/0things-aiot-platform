package data

import (
	"context"

	"telemetry-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type eventRepo struct {
	data   *Data
	logger *log.Helper
}

// NewEventRepo creates a new event repository.
func NewEventRepo(data *Data, logger log.Logger) biz.EventRepo {
	return &eventRepo{
		data:   data,
		logger: log.NewHelper(logger),
	}
}

// Save stores an event.
func (r *eventRepo) Save(ctx context.Context, event *biz.Event) error {
	r.logger.Infof("Saving event: type=%s, product_key=%s, device_key=%s",
		event.Type, event.ProductKey, event.DeviceKey)

	// TODO: Implement actual database storage
	r.logger.Debugf("Event data: %+v", event)

	return nil
}
