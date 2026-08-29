package server

import (
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/pkg/log"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
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
		&model.Organization{},
		&model.OrganizationUser{},
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
	// Rename OTA batch/device-status tables to the ota_-prefixed names
	// while preserving existing data. Only runs when the old table exists.
	for _, rename := range []struct{ from, to string }{
		{"upgrade_batches", "ota_upgrade_batches"},
		{"device_upgrade_status", "ota_device_upgrade_status"},
	} {
		if m.deviceDB.Migrator().HasTable(rename.from) {
			if err := m.deviceDB.Migrator().RenameTable(rename.from, rename.to); err != nil {
				return fmt.Errorf("rename %s -> %s: %w", rename.from, rename.to, err)
			}
		}
	}
	// Drop the unused batch_type column from ota_upgrade_batches. Guarded so
	// re-running migration is a no-op once the column is gone.
	if m.deviceDB.Migrator().HasColumn("ota_upgrade_batches", "batch_type") {
		if err := m.deviceDB.Exec("ALTER TABLE ota_upgrade_batches DROP COLUMN batch_type").Error; err != nil {
			return fmt.Errorf("drop batch_type column: %w", err)
		}
	}
	// Rename last_status_change_time -> last_status_change_ts on the device
	// upgrade status table while preserving existing data. Guarded so
	// re-running migration is a no-op once the column is renamed.
	if m.deviceDB.Migrator().HasColumn("ota_device_upgrade_status", "last_status_change_time") {
		if err := m.deviceDB.Migrator().RenameColumn(
			"ota_device_upgrade_status", "last_status_change_time", "last_status_change_ts",
		); err != nil {
			return fmt.Errorf("rename last_status_change_time -> last_status_change_ts: %w", err)
		}
	}
	if err := m.deviceDB.AutoMigrate(
		&model.Product{},
		&model.Category{},
		&model.ProductProtocol{},
		&model.ProductTSL{},
		&model.ProductMessageParser{},
		&model.Device{},
		&model.DeviceEndpoint{},
		&model.DeviceGroup{},
		&model.DeviceGroupMember{},
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
	if err := m.seedDefaultCategories(); err != nil {
		return fmt.Errorf("seed product categories: %w", err)
	}
	for _, column := range []string{"config_schema", "support_uplink", "support_downlink", "support_ota", "enabled"} {
		if m.deviceDB.Migrator().HasColumn("product_protocols", column) {
			if err := m.deviceDB.Exec("ALTER TABLE product_protocols DROP COLUMN " + column).Error; err != nil {
				return fmt.Errorf("drop product_protocols.%s: %w", column, err)
			}
		}
	}
	// 将旧产品上的单值接入协议迁移为产品协议能力。
	if err := m.backfillProductProtocols(); err != nil {
		return fmt.Errorf("backfill product protocols: %w", err)
	}
	for _, table := range []string{"products", "devices"} {
		if err := m.deviceDB.Exec("UPDATE " + table + " SET organization_id = 1 WHERE organization_id IS NULL OR organization_id = 0").Error; err != nil {
			return fmt.Errorf("backfill organization_id for %s: %w", table, err)
		}
	}
	if !m.deviceDB.Migrator().HasTable(&model.DeviceEvent{}) {
		return fmt.Errorf("device_events migration did not create the table")
	}
	// Backfill the uuid column for any existing OTA packages that lack one.
	if err := m.backfillPackageUUIDs(); err != nil {
		return fmt.Errorf("backfill package uuid: %w", err)
	}
	m.log.Info("AutoMigrate success")
	fmt.Println("migration complete: device_events table is ready")
	os.Exit(0)
	return nil
}

// seedDefaultCategories keeps existing product creation usable after the category tree is introduced.
func (m *MigrateServer) seedDefaultCategories() error {
	var count int64
	if err := m.deviceDB.Model(&model.Category{}).Count(&count).Error; err != nil || count > 0 {
		if err != nil {
			return err
		}
		return nil
	}
	items := []string{"传感器", "执行器", "网关", "控制器", "显示设备", "摄像头", "其他"}
	for i, name := range items {
		if err := m.deviceDB.Create(&model.Category{Name: name, Sort: i, Enabled: true}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (m *MigrateServer) backfillProductProtocols() error {
	var products []model.Product
	if err := m.deviceDB.Where("access_protocol IS NOT NULL AND access_protocol <> ?", "").Find(&products).Error; err != nil {
		return err
	}
	for _, product := range products {
		var count int64
		if err := m.deviceDB.Model(&model.ProductProtocol{}).Where("product_id = ?", product.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		transports := []string{product.AccessProtocol}
		if product.AccessProtocol == "default" {
			transports = []string{"http", "mqtt"}
		}
		for _, transport := range transports {
			if err := m.deviceDB.Create(&model.ProductProtocol{
				ProductID: product.ID, TransportProtocol: transport,
				ApplicationProtocol: "json",
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}

// backfillPackageUUIDs assigns a generated uuid to every OTA package that does
// not yet have one. Safe to re-run: rows with a non-empty uuid are skipped.
func (m *MigrateServer) backfillPackageUUIDs() error {
	var packages []model.OTAPackage
	if err := m.deviceDB.Where("uuid IS NULL OR uuid = ?", "").Find(&packages).Error; err != nil {
		return err
	}
	for _, pkg := range packages {
		if err := m.deviceDB.Model(&model.OTAPackage{}).Where("id = ?", pkg.ID).
			Update("uuid", uuid.NewString()).Error; err != nil {
			return err
		}
	}
	return nil
}
