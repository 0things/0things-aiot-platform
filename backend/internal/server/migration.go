package server

import (
	"context"

	"aiot-backend/internal/model"
	"aiot-backend/pkg/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MigrateServer struct {
	db  *gorm.DB
	log *log.Logger
}

func NewMigrateServer(db *gorm.DB, log *log.Logger) *MigrateServer {
	return &MigrateServer{
		db:  db,
		log: log,
	}
}

// Start executes GORM AutoMigrate for all system and IoT models.
func (m *MigrateServer) Start(ctx context.Context) error {
	m.log.Info("Starting GORM AutoMigrate...")

	if err := m.db.WithContext(ctx).AutoMigrate(
		&model.User{},
		&model.Organization{},
		&model.OrganizationUser{},
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
		&model.DeviceServiceInvocation{},
		&model.DevicePushRecord{},
		&model.OTAPackage{},
		&model.UpgradeBatch{},
		&model.DeviceUpgradeStatus{},
		&model.SceneLinkage{},
		&model.SceneLinkageDetail{},
	); err != nil {
		m.log.Error("AutoMigrate failed", zap.Error(err))
		return err
	}

	m.log.Info("GORM AutoMigrate completed successfully")
	return nil
}

func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("MigrateServer stopped")
	return nil
}
