package testutil

import (
	"context"
	"testing"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/service"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type testKafkaService struct{}

func (testKafkaService) Produce(context.Context, string, []byte, []byte) error              { return nil }
func (testKafkaService) ProduceJSON(context.Context, string, string, any) error             { return nil }
func (testKafkaService) ProduceAsync(context.Context, string, []byte, []byte, func(error))  {}
func (testKafkaService) ProduceJSONAsync(context.Context, string, string, any, func(error)) {}
func (testKafkaService) Flush(context.Context) error                                        { return nil }
func (testKafkaService) Close()                                                             {}

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.Category{},
		&model.ProductProtocol{},
		&model.ProductTSL{},
		&model.Device{},
		&model.DeviceState{},
		&model.DeviceTag{},
		&model.DeviceShadow{},
		&model.DeviceShadowHistory{},
		&model.DeviceEvent{},
		&model.DeviceServiceInvocation{},
		&model.DevicePushRecord{},
		&model.OTAPackage{},
		&model.UpgradeBatch{},
		&model.DeviceUpgradeStatus{},
		&model.SceneLinkage{},
		&model.SceneLinkageDetail{},
	)
	require.NoError(t, err)

	return db
}

func SeedTestData(t *testing.T, db *gorm.DB) {
	t.Helper()
	product := &model.Product{ID: 1, ProductKey: "P001", Name: "Test Product", OrganizationID: 1}
	db.Create(product)

	device := &model.Device{ID: 1, DeviceKey: "D001", Name: "Test Device", ProductID: 1, OrganizationID: 1, Enabled: true}
	db.Create(device)

	db.Create(&model.DeviceState{ID: 1, DeviceKey: "D001", State: "online"})
}

func NewTestRepositories(db *gorm.DB) (*repository.DeviceRepository, *repository.ProductRepository, *repository.DeviceTagRepository, *repository.DeviceShadowRepository, *repository.PushRecordRepository) {
	deviceRepo := repository.NewDeviceRepository(db, nil)
	productRepo := repository.NewProductRepository(db)
	tagRepo := repository.NewDeviceTagRepository(db)
	shadowRepo := repository.NewDeviceShadowRepository(db)
	pushRepo := repository.NewPushRecordRepository(db)

	return deviceRepo, productRepo, tagRepo, shadowRepo, pushRepo
}

func NewTestDeviceService(db *gorm.DB) *service.DeviceService {
	deviceRepo, productRepo, tagRepo, shadowRepo, pushRepo := NewTestRepositories(db)
	return service.NewDeviceService(deviceRepo, productRepo, tagRepo, shadowRepo, pushRepo)
}

func NewTestProductService(db *gorm.DB) *service.ProductService {
	productRepo := repository.NewProductRepository(db)
	return service.NewProductService(productRepo)
}

func NewTestOTATotalService(db *gorm.DB) *service.OTAService {
	otaRepo := repository.NewOTARepository(db)
	productRepo := repository.NewProductRepository(db)
	deviceRepo := repository.NewDeviceRepository(db, nil)
	return service.NewOTAService(otaRepo, productRepo, deviceRepo, testKafkaService{})
}

func NewTestSceneLinkageService(db *gorm.DB) *service.SceneLinkageService {
	sceneRepo := repository.NewSceneLinkageRepository(db)
	return service.NewSceneLinkageService(sceneRepo)
}

func NewTestDeviceEventService(db *gorm.DB) *service.DeviceEventService {
	eventRepo := repository.NewDeviceEventRepository(db)
	deviceRepo := repository.NewDeviceRepository(db, nil)
	return service.NewDeviceEventService(eventRepo, deviceRepo)
}

func NewTestThingModelDataService(db *gorm.DB) *service.ThingModelDataService {
	invocations := repository.NewDeviceServiceInvocationRepository(db)
	deviceRepo := repository.NewDeviceRepository(db, nil)
	tslRepo := repository.NewProductTSLRepository(db)
	return service.NewThingModelDataService(invocations, deviceRepo, tslRepo, nil)
}

func NewTestProductTSLService(db *gorm.DB) *service.ProductTSLService {
	productRepo := repository.NewProductRepository(db)
	tslRepo := repository.NewProductTSLRepository(db)
	return service.NewProductTSLService(productRepo, tslRepo)
}

func ContextWithTenant(ctx context.Context, organizationID int64) context.Context {
	return context.WithValue(ctx, "organization_id", organizationID)
}
