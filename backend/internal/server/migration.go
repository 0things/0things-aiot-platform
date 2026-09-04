package server

import (
	"context"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/pkg/log"
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

// Start executes GORM AutoMigrate for all system and IoT models.
func (m *MigrateServer) Start(ctx context.Context) error {
	m.log.Info("Starting GORM AutoMigrate...")

	// 1. User & IAM DB Schema
	if err := m.db.WithContext(ctx).AutoMigrate(
		&model.User{},
		&model.Organization{},
		&model.OrganizationUser{},
	); err != nil {
		m.log.Error("User DB AutoMigrate failed", zap.Error(err))
		return err
	}

	// 2. IoT & Device Management DB Schema
	if err := m.deviceDB.WithContext(ctx).AutoMigrate(
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
		m.log.Error("Device DB AutoMigrate failed", zap.Error(err))
		return err
	}

	// 3. Seed minimal default metadata if database is empty
	if err := m.seedDefaults(ctx); err != nil {
		m.log.Warn("Failed to seed default metadata", zap.Error(err))
	}

	m.log.Info("GORM AutoMigrate completed successfully")
	return nil
}

// seedDefaults initializes default categories when empty.
func (m *MigrateServer) seedDefaults(ctx context.Context) error {
	var count int64
	if err := m.deviceDB.WithContext(ctx).Model(&model.Category{}).Count(&count).Error; err == nil && count == 0 {
		items := []string{"传感器", "执行器", "网关", "控制器", "显示设备", "摄像头", "其他"}
		for i, name := range items {
			_ = m.deviceDB.WithContext(ctx).Create(&model.Category{Name: name, Sort: i, Enabled: true}).Error
		}
	}
	return nil
}

func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("MigrateServer stopped")
	return nil
}
