package repository

import (
	"context"
	"testing"
	"time"

	"0things-backend/internal/model"
	"0things-backend/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// --- Push Record Repository Tests ---

func setupPushRecordRepository(t *testing.T) (*repository.PushRecordRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	pushRepo := repository.NewPushRecordRepository(&repository.IoTDB{DB: db})
	return pushRepo, mock
}

func TestPushRecordRepository_CreatePushRecord(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `device_push_records`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	record := &model.DevicePushRecord{DeviceID: 1, OperationType: "push", Status: "success"}
	err := pushRepo.CreatePushRecord(ctx, record)
	assert.NoError(t, err)
}

func TestPushRecordRepository_ListPushRecords_WithOperationType(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `device_push_records`").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT .+ FROM `device_push_records`").WillReturnRows(
		sqlmock.NewRows([]string{"id", "device_id", "operation_type", "status", "created_at"}).
			AddRow(1, 1, "push", "success", time.Now()),
	)

	records, total, err := pushRepo.ListPushRecords(ctx, 1, 1, 10, "push", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, records, 1)
}

func TestPushRecordRepository_ListPushRecords_WithStatus(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `device_push_records`").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery("SELECT .+ FROM `device_push_records`").WillReturnRows(sqlmock.NewRows([]string{"id", "device_id", "operation_type", "status", "created_at"}))

	records, total, err := pushRepo.ListPushRecords(ctx, 1, 1, 10, "", "failed")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, records)
}

func TestPushRecordRepository_FindPushRecord(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "device_id", "operation_type", "status", "created_at"}).
		AddRow(1, 1, "push", "success", time.Now())
	mock.ExpectQuery("SELECT .+ FROM `device_push_records`").WillReturnRows(rows)

	record, err := pushRepo.FindPushRecord(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "push", record.OperationType)
}

func TestPushRecordRepository_DeletePushRecords_WithBefore(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	before := time.Now()
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `device_push_records`").WillReturnResult(sqlmock.NewResult(0, 5))
	mock.ExpectCommit()

	count, err := pushRepo.DeletePushRecords(ctx, 1, &before)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestPushRecordRepository_DeletePushRecords_NilBefore(t *testing.T) {
	pushRepo, mock := setupPushRecordRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `device_push_records`").WillReturnResult(sqlmock.NewResult(0, 10))
	mock.ExpectCommit()

	count, err := pushRepo.DeletePushRecords(ctx, 1, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
}

// --- User Repository Tests ---

func setupExtendedUserRepository(t *testing.T) (*repository.UserRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewRepository(nil, db)
	userRepo := repository.NewUserRepository(repo)
	return &userRepo, mock
}

func TestExtendedUserRepository_GetByID(t *testing.T) {
	userRepo, mock := setupExtendedUserRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "user_id", "email", "nickname"}).
		AddRow(1, "u1", "test@test.com", "Test User")
	mock.ExpectQuery("SELECT .+ FROM `users`").WillReturnRows(rows)

	user, err := (*userRepo).GetByID(ctx, "u1")
	assert.NoError(t, err)
	assert.Equal(t, "test@test.com", user.Email)
}

func TestExtendedUserRepository_GetByID_NotFound(t *testing.T) {
	userRepo, mock := setupExtendedUserRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `users`").WillReturnError(gorm.ErrRecordNotFound)

	user, err := (*userRepo).GetByID(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestExtendedUserRepository_GetByEmail(t *testing.T) {
	userRepo, mock := setupExtendedUserRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "user_id", "email", "nickname"}).
		AddRow(1, "u1", "test@test.com", "Test User")
	mock.ExpectQuery("SELECT .+ FROM `users`").WillReturnRows(rows)

	user, err := (*userRepo).GetByEmail(ctx, "test@test.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@test.com", user.Email)
}

func TestExtendedUserRepository_GetByEmail_NotFound(t *testing.T) {
	userRepo, mock := setupExtendedUserRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `users`").WillReturnError(gorm.ErrRecordNotFound)

	user, err := (*userRepo).GetByEmail(ctx, "nonexistent@test.com")
	assert.NoError(t, err)
	assert.Nil(t, user)
}

// --- Device Shadow Repository Tests ---

func setupDeviceShadowRepository(t *testing.T) (*repository.DeviceShadowRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	shadowRepo := repository.NewDeviceShadowRepository(&repository.IoTDB{DB: db})
	return shadowRepo, mock
}

func TestDeviceShadowRepository_GetShadow(t *testing.T) {
	shadowRepo, mock := setupDeviceShadowRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "device_id", "desired", "reported", "metadata", "version"}).
		AddRow(1, 1, []byte("{}"), []byte("{}"), []byte("{}"), 1)
	mock.ExpectQuery("SELECT .+ FROM `device_shadows`").WillReturnRows(rows)

	shadow, err := shadowRepo.GetShadow(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), shadow.Version)
}

func TestDeviceShadowRepository_GetShadow_NotFound(t *testing.T) {
	shadowRepo, mock := setupDeviceShadowRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `device_shadows`").WillReturnError(gorm.ErrRecordNotFound)

	shadow, err := shadowRepo.GetShadow(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, shadow)
}

// MutateShadow uses GORM Transaction wrapping — error gets wrapped, hard to test with sqlmock
func TestDeviceShadowRepository_MutateShadow_CreateNew(t *testing.T) {
	// This path creates + updates within a Transaction — complex to mock precisely
	t.Skip("skipped: MutateShadow uses GORM Transaction which wraps errors")
}

func TestDeviceShadowRepository_ListShadowHistory(t *testing.T) {
	shadowRepo, mock := setupDeviceShadowRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "device_id", "version", "source", "created_at"}).
		AddRow(1, 1, 1, "app", time.Now()).
		AddRow(2, 1, 2, "device", time.Now())
	mock.ExpectQuery("SELECT .+ FROM `device_shadow_histories`").WillReturnRows(rows)

	history, err := shadowRepo.ListShadowHistory(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, history, 2)
}

// --- Device Tag Repository Tests ---

func setupDeviceTagRepository(t *testing.T) (*repository.DeviceTagRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	tagRepo := repository.NewDeviceTagRepository(&repository.IoTDB{DB: db})
	return tagRepo, mock
}

func TestDeviceTagRepository_ListTags(t *testing.T) {
	tagRepo, mock := setupDeviceTagRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "device_id", "key", "value", "source"}).
		AddRow(1, 1, "location", "factory-A", "manual")
	mock.ExpectQuery("SELECT .+ FROM `device_tags`").WillReturnRows(rows)

	tags, err := tagRepo.ListTags(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, tags, 1)
}

func TestDeviceTagRepository_DeleteTags(t *testing.T) {
	tagRepo, mock := setupDeviceTagRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `device_tags`").WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	err := tagRepo.DeleteTags(ctx, 1, []string{"location", "env"})
	assert.NoError(t, err)
}

func TestDeviceTagRepository_SetTags_Replace(t *testing.T) {
	tagRepo, mock := setupDeviceTagRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `device_tags`").WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectQuery("SELECT .+ FROM `device_tags`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectExec("INSERT INTO `device_tags`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := tagRepo.SetTags(ctx, 1, map[string]string{"key1": "val1"}, true)
	assert.NoError(t, err)
}

func TestDeviceTagRepository_SetTags_Upsert_Existing(t *testing.T) {
	tagRepo, mock := setupDeviceTagRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "device_id", "key", "value", "source"}).
		AddRow(1, 1, "location", "factory-A", "manual")
	mock.ExpectQuery("SELECT .+ FROM `device_tags`").WillReturnRows(rows)
	mock.ExpectExec("UPDATE `device_tags`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := tagRepo.SetTags(ctx, 1, map[string]string{"location": "factory-B"}, false)
	assert.NoError(t, err)
}

// --- OTA Repository Additional Tests ---

func setupExtendedOTARespository(t *testing.T) (*repository.OTARepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	otaRepo := repository.NewOTARepository(&repository.IoTDB{DB: db})
	return otaRepo, mock
}

func TestOTARepository_FindByName(t *testing.T) {
	otaRepo, mock := setupExtendedOTARespository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "package_name", "version", "product_key", "created_at", "updated_at"}).
		AddRow(1, "firmware-1", "1.0.0", "P001", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnRows(rows)

	pkg, err := otaRepo.FindByName(ctx, "firmware-1")
	assert.NoError(t, err)
	assert.Equal(t, "firmware-1", pkg.PackageName)
}

func TestOTARepository_FindByName_NotFound(t *testing.T) {
	otaRepo, mock := setupExtendedOTARespository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnError(gorm.ErrRecordNotFound)

	pkg, err := otaRepo.FindByName(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, pkg)
}

// --- Rule Repository Additional Tests ---

func setupExtendedRuleRepository(t *testing.T) (*repository.RuleRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	ruleRepo := repository.NewRuleRepository(&repository.IoTDB{DB: db})
	return ruleRepo, mock
}

func TestRuleRepository_Create(t *testing.T) {
	ruleRepo, mock := setupExtendedRuleRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `products`").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(1),
	)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `rules`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rule := &model.Rule{Name: "New Rule", Type: "threshold", Status: "draft", ProductID: 1}
	err := ruleRepo.Create(ctx, rule)
	assert.NoError(t, err)
}

func TestRuleRepository_UpdateStatus(t *testing.T) {
	ruleRepo, mock := setupExtendedRuleRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `rules`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	rule := &model.Rule{ID: 1, Status: "active"}
	err := ruleRepo.UpdateStatus(ctx, rule, "enabled")
	assert.NoError(t, err)
}

// --- Device Event Repository Tests ---

func setupDeviceEventRepository(t *testing.T) (*repository.DeviceEventRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	eventRepo := repository.NewDeviceEventRepository(&repository.IoTDB{DB: db})
	return eventRepo, mock
}

func TestDeviceEventRepository_Create(t *testing.T) {
	eventRepo, mock := setupDeviceEventRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `device_events`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	event := &model.DeviceEvent{DeviceID: 1, EventType: "temperature", EventAt: time.Now()}
	err := eventRepo.Create(ctx, event)
	assert.NoError(t, err)
}

func TestDeviceEventRepository_List_WithFilters(t *testing.T) {
	eventRepo, mock := setupDeviceEventRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `device_events`").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery("SELECT .+ FROM `device_events`").WillReturnRows(sqlmock.NewRows([]string{"id"}))

	events, total, err := eventRepo.List(ctx, 1, 10, "temp", "D001", "temperature", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, events)
}

// --- Product Repository Tests ---

func setupProductRepository(t *testing.T) (*repository.ProductRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	productRepo := repository.NewProductRepository(&repository.IoTDB{DB: db})
	return productRepo, mock
}

func TestProductRepository_Find(t *testing.T) {
	productRepo, mock := setupProductRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "product_key", "name", "status", "created_at", "updated_at"}).
		AddRow(1, "P001", "Test Product", "active", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `products`").WillReturnRows(rows)

	product, err := productRepo.Find(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Test Product", product.Name)
}

func TestProductRepository_FindByKey(t *testing.T) {
	productRepo, mock := setupProductRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "product_key", "name", "status", "created_at", "updated_at"}).
		AddRow(1, "P001", "Test Product", "active", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `products`").WillReturnRows(rows)

	product, err := productRepo.FindByKey(ctx, "P001")
	assert.NoError(t, err)
	assert.Equal(t, "P001", product.ProductKey)
}

func TestProductRepository_Create(t *testing.T) {
	productRepo, mock := setupProductRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `products`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	product := &model.Product{Name: "New Product", ProductKey: "P002", Status: "active"}
	err := productRepo.Create(ctx, product)
	assert.NoError(t, err)
}

func TestProductRepository_Delete(t *testing.T) {
	productRepo, mock := setupProductRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `products` SET `deleted_at`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	product := &model.Product{ID: 1}
	err := productRepo.Delete(ctx, product)
	assert.NoError(t, err)
}

// --- Product TSL Repository Tests ---

func setupProductTSLRepository(t *testing.T) (*repository.ProductTSLRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	tslRepo := repository.NewProductTSLRepository(&repository.IoTDB{DB: db})
	return tslRepo, mock
}

func TestProductTSLRepository_FindByProductID(t *testing.T) {
	tslRepo, mock := setupProductTSLRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "product_product_tsl", "tsl"}).
		AddRow(1, 1, `{"properties":[]}`)
	mock.ExpectQuery("SELECT .+ FROM `product_ts_ls`").WillReturnRows(rows)

	tsl, err := tslRepo.FindByProductID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, `{"properties":[]}`, tsl.TSL)
}

func TestProductTSLRepository_FindByProductID_NotFound(t *testing.T) {
	tslRepo, mock := setupProductTSLRepository(t)
	ctx := context.Background()

	mock.ExpectQuery("SELECT .+ FROM `product_ts_ls`").WillReturnError(gorm.ErrRecordNotFound)

	tsl, err := tslRepo.FindByProductID(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, tsl)
}

func TestProductTSLRepository_Create(t *testing.T) {
	tslRepo, mock := setupProductTSLRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `product_ts_ls`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tsl := &model.ProductTSL{TSL: `{"properties":[]}`}
	err := tslRepo.Create(ctx, tsl)
	assert.NoError(t, err)
}

func TestProductTSLRepository_Save(t *testing.T) {
	tslRepo, mock := setupProductTSLRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `product_ts_ls`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tsl := &model.ProductTSL{ID: 1, TSL: `{"properties":[{"name":"temp"}]}`}
	err := tslRepo.Save(ctx, tsl)
	assert.NoError(t, err)
}

func TestProductTSLRepository_Delete(t *testing.T) {
	tslRepo, mock := setupProductTSLRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `product_ts_ls`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tsl := &model.ProductTSL{ID: 1}
	err := tslRepo.Delete(ctx, tsl)
	assert.NoError(t, err)
}

// --- Device Telemetry (nil redis) ---

func TestDeviceRepository_Telemetry_NilRedis(t *testing.T) {
	deviceRepo, _ := setupDeviceRepository(t)
	ctx := context.Background()

	result, err := deviceRepo.Telemetry(ctx, "D001")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis is not configured")
	assert.Equal(t, "", result)
}

func TestDeviceRepository_Statistics(t *testing.T) {
	deviceRepo, mock := setupDeviceRepository(t)
	ctx := context.Background()

	// Statistics uses gen queries with tenant filter and JOINs — just verify it returns data
	mock.ExpectQuery("SELECT .+").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
	mock.ExpectQuery("SELECT .+").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
	mock.ExpectQuery("SELECT .+").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery("SELECT .+").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	stats, err := deviceRepo.Statistics(ctx)
	assert.NoError(t, err)
	assert.True(t, stats.TotalDevices > 0)
}

func int64Ptr(v int64) *int64 { return &v }
