package server

import (
	"0things-backend/internal/model"
	"0things-backend/internal/repository"
	"0things-backend/pkg/log"
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gorm.io/gorm"
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
	if err := m.deviceDB.AutoMigrate(
		&model.Product{},
		&model.ProductTSL{},
		&model.Device{},
		&model.DeviceState{},
		&model.DeviceShadow{},
		&model.DeviceTag{},
		&model.DeviceShadowHistory{},
		&model.DeviceEvent{},
		&model.DevicePushRecord{},
		&model.OTAPackage{},
		&model.UpgradeBatch{},
		&model.DeviceUpgradeStatus{},
		&model.SceneLinkage{},
		&model.SceneLinkageDetail{},
		&model.Alert{},
	); err != nil {
		return err
	}
	for _, table := range []string{"products", "devices"} {
		if err := m.deviceDB.Exec("UPDATE " + table + " SET tenant_id = 1 WHERE tenant_id IS NULL OR tenant_id = 0").Error; err != nil {
			return fmt.Errorf("backfill tenant_id for %s: %w", table, err)
		}
	}
	if !m.deviceDB.Migrator().HasTable(&model.DeviceEvent{}) {
		return fmt.Errorf("device_events migration did not create the table")
	}
	m.log.Info("AutoMigrate success")
	fmt.Println("migration complete: device_events table is ready")
	os.Exit(0)
	return nil
}
func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}
