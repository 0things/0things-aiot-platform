package server

import (
	"0things-backend/internal/model"
	"0things-backend/internal/repository"
	"0things-backend/pkg/log"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

type MigrateServer struct {
	db       *gorm.DB
	deviceDB *repository.IoTDB
	log      *log.Logger
}

func NewMigrateServer(db *gorm.DB, deviceDB *repository.IoTDB, log *log.Logger) *MigrateServer {
	return &MigrateServer{
		db:       db,
		deviceDB: deviceDB,
		log:      log,
	}
}
func (m *MigrateServer) Start(ctx context.Context) error {
	if err := m.db.AutoMigrate(
		&model.User{},
	); err != nil {
		m.log.Error("user migrate error", zap.Error(err))
		return err
	}
	for _, table := range []string{"products", "devices", "device_states", "product_ts_ls"} {
		if !m.deviceDB.Migrator().HasTable(table) {
			return fmt.Errorf("legacy device table %q is missing; refusing to create legacy data tables", table)
		}
	}
	// These three tables were designed in device-service but were not wired
	// there. AutoMigrate only creates missing tables; it does not touch the
	// pre-existing legacy product/device/OTA/rule tables validated above.
	if err := m.deviceDB.AutoMigrate(&model.DeviceShadow{}, &model.DeviceTag{}, &model.DeviceShadowHistory{}); err != nil {
		return err
	}
	m.log.Info("AutoMigrate success")
	os.Exit(0)
	return nil
}
func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}
