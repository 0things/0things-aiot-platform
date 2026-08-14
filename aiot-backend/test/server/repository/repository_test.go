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

func setupAlertRepository(t *testing.T) (*repository.AlertRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	alertRepo := repository.NewAlertRepository(&repository.IoTDB{DB: db})
	return alertRepo, mock
}

func TestAlertRepository_List(t *testing.T) {
	alertRepo, mock := setupAlertRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "summary", "status", "severity", "device_key", "created_at", "updated_at"}).
		AddRow(1, "High Temperature", "active", "critical", "D001", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `alerts`").WillReturnRows(rows)

	alerts, total, err := alertRepo.List(ctx, 1, 10, "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, alerts, 1)
}

func TestAlertRepository_Find(t *testing.T) {
	alertRepo, mock := setupAlertRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "summary", "status", "severity", "device_key", "created_at", "updated_at"}).
		AddRow(1, "High Temperature", "active", "critical", "D001", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `alerts`").WillReturnRows(rows)

	alert, err := alertRepo.Find(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "High Temperature", alert.Summary)
}

func TestAlertRepository_UpdateStatus(t *testing.T) {
	alertRepo, mock := setupAlertRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `alerts`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	alert := &model.Alert{ID: 1, Summary: "High Temperature"}
	err := alertRepo.UpdateStatus(ctx, alert, "resolved", time.Now())
	assert.NoError(t, err)
}

func setupDeviceRepository(t *testing.T) (*repository.DeviceRepository, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	deviceRepo := repository.NewDeviceRepository(&repository.IoTDB{DB: db}, &repository.IoTRedis{Client: nil})
	return deviceRepo, mock
}

func TestDeviceRepository_DB(t *testing.T) {
	deviceRepo, _ := setupDeviceRepository(t)
	ctx := context.Background()

	db := deviceRepo.DB(ctx)
	assert.NotNil(t, db)
}

func TestDeviceRepository_Create(t *testing.T) {
	deviceRepo, mock := setupDeviceRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `devices`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO `device_states`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	device := &model.Device{DeviceKey: "D002", Name: "New Device", ProductID: 1, TenantID: 1}
	err := deviceRepo.Create(ctx, device)
	assert.NoError(t, err)
}

func TestDeviceRepository_Delete(t *testing.T) {
	deviceRepo, mock := setupDeviceRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `devices` SET `deleted_at`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	device := &model.Device{ID: 1}
	err := deviceRepo.Delete(ctx, device)
	assert.NoError(t, err)
}

func TestDeviceRepository_Restore(t *testing.T) {
	deviceRepo, mock := setupDeviceRepository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `devices`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := deviceRepo.Restore(ctx, 1)
	assert.NoError(t, err)
}

// Skip TestDeviceRepository_SaveEnabled - GORM Save uses upsert which is hard to mock precisely with sqlmock

func setupOTARespository(t *testing.T) (*repository.OTARepository, sqlmock.Sqlmock) {
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

func TestOTARepository_List(t *testing.T) {
	otaRepo, mock := setupOTARespository(t)
	ctx := context.Background()

	// Count query
	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	// Find query
	rows := sqlmock.NewRows([]string{"id", "package_name", "version", "product_key", "created_at", "updated_at"}).
		AddRow(1, "firmware-1", "1.0.0", "P001", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnRows(rows)

	packages, total, err := otaRepo.List(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, packages, 1)
}

func TestOTARepository_Find(t *testing.T) {
	otaRepo, mock := setupOTARespository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "package_name", "version", "product_key", "created_at", "updated_at"}).
		AddRow(1, "firmware-1", "1.0.0", "P001", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnRows(rows)

	pkg, err := otaRepo.Find(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "firmware-1", pkg.PackageName)
}

func TestOTARepository_Create(t *testing.T) {
	otaRepo, mock := setupOTARespository(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `ota_packages`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	pkg := &model.OTAPackage{PackageName: "firmware-1", ProductKey: "P001", Version: "1.0.0"}
	err := otaRepo.Create(ctx, pkg)
	assert.NoError(t, err)
}

func TestOTARepository_Delete(t *testing.T) {
	otaRepo, mock := setupOTARespository(t)
	ctx := context.Background()

	// Find is called first inside Delete
	rows := sqlmock.NewRows([]string{"id", "package_name", "version", "product_key", "created_at", "updated_at"}).
		AddRow(1, "firmware-1", "1.0.0", "P001", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `ota_packages`").WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `ota_packages` SET `deleted_at`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := otaRepo.Delete(ctx, 1)
	assert.NoError(t, err)
}

func setupRuleRepository(t *testing.T) (*repository.RuleRepository, sqlmock.Sqlmock) {
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

func TestRuleRepository_Find(t *testing.T) {
	ruleRepo, mock := setupRuleRepository(t)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "product_id", "status", "created_at", "updated_at"}).
		AddRow(1, "Temperature Alert", "Alert", 1, "active", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `rules`").WillReturnRows(rows)

	rule, err := ruleRepo.Find(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Temperature Alert", rule.Name)
}

func TestRuleRepository_Delete(t *testing.T) {
	ruleRepo, mock := setupRuleRepository(t)
	ctx := context.Background()

	// Find is called first inside Delete
	rows := sqlmock.NewRows([]string{"id", "name", "description", "product_id", "status", "created_at", "updated_at"}).
		AddRow(1, "Temperature Alert", "Alert", 1, "active", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `rules`").WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `rules`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := ruleRepo.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestRuleRepository_CreateExecution(t *testing.T) {
	ruleRepo, mock := setupRuleRepository(t)
	ctx := context.Background()

	// First expect the Find query
	rows := sqlmock.NewRows([]string{"id", "name", "description", "product_id", "status", "created_at", "updated_at"}).
		AddRow(1, "Temperature Alert", "Alert", 1, "active", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .+ FROM `rules`").WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `rule_executions`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	execution := &model.RuleExecution{RuleID: 1, RuleName: "Temperature Alert", Status: "success"}
	err := ruleRepo.CreateExecution(ctx, execution)
	assert.NoError(t, err)
}
