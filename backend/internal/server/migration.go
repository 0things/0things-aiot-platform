package server

import (
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/pkg/log"
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
	// Rename the legacy tenant_id column to organization_id while preserving data.
	for _, table := range []string{"products", "devices", "product_message_parsers", "ota_packages", "scene_linkage"} {
		if m.deviceDB.Migrator().HasColumn(table, "tenant_id") {
			if err := m.deviceDB.Migrator().RenameColumn(table, "tenant_id", "organization_id"); err != nil {
				return fmt.Errorf("rename tenant_id -> organization_id for %s: %w", table, err)
			}
		}
	}
	if err := m.deviceDB.AutoMigrate(
		&model.Product{},
		&model.ProductTSL{},
		&model.ProductMessageParser{},
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
	); err != nil {
		return err
	}
	for _, table := range []string{"products", "devices"} {
		if err := m.deviceDB.Exec("UPDATE " + table + " SET organization_id = 1 WHERE organization_id IS NULL OR organization_id = 0").Error; err != nil {
			return fmt.Errorf("backfill organization_id for %s: %w", table, err)
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
