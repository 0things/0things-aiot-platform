package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"0things-backend/internal/model"
	"0things-backend/internal/repository"
	"0things-backend/internal/service"
	"0things-backend/internal/tenant"
	"0things-backend/test/server/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func ctx2() context.Context {
	return tenant.WithTenant(context.Background(), 1)
}

func TestIntegrationDeviceService_CreateDevice_InvalidMetadata(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	device := &model.Device{Name: "Bad Meta", ProductID: 1, TenantID: 1, Metadata: json.RawMessage(`"not-valid-json`)}
	_, err := svc.CreateDevice(ctx2(), device)
	assert.Error(t, err)
}

func TestIntegrationDeviceService_CreateDevice_WithValidMetadata(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	device := &model.Device{Name: "Good Meta", ProductID: 1, TenantID: 1, Metadata: json.RawMessage(`{"key":"value"}`)}
	result, err := svc.CreateDevice(ctx2(), device)
	require.NoError(t, err)
	assert.NotZero(t, result.ID)
}

func TestIntegrationDeviceService_CreateDevice_WithStringMetadata(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	device := &model.Device{Name: "String Meta", ProductID: 1, TenantID: 1, Metadata: json.RawMessage(`"{\"key\":\"value\"}"`)}
	result, err := svc.CreateDevice(ctx2(), device)
	require.NoError(t, err)
	assert.NotZero(t, result.ID)
}

func TestIntegrationDeviceService_CreateDevice_WithCustomKey(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	device := &model.Device{Name: "Custom Key", ProductID: 1, TenantID: 1, DeviceKey: "CUSTOM001"}
	result, err := svc.CreateDevice(ctx2(), device)
	require.NoError(t, err)
	assert.Equal(t, "CUSTOM001", result.DeviceKey)
}

func TestIntegrationDeviceService_CreateDevice_WrongProduct(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	device := &model.Device{Name: "Bad Product", ProductID: 999, TenantID: 1}
	_, err := svc.CreateDevice(ctx2(), device)
	assert.Error(t, err)
}

func TestIntegrationDeviceService_Device_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.Device(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationDeviceService_DeviceByKey_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.DeviceByKey(ctx2(), "NONEXIST")
	assert.Error(t, err)
}

func TestIntegrationDeviceService_UpdateDevice_InvalidTransition(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	// Device is "online", trying to go to "inactive" should fail
	_, err := svc.UpdateDevice(ctx2(), 1, "", "inactive", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status transition")
}

func TestIntegrationDeviceService_UpdateDevice_ValidTransition(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	// Device is "online", can go to "offline"
	device, err := svc.UpdateDevice(ctx2(), 1, "", "offline", nil)
	require.NoError(t, err)
	assert.Equal(t, "offline", device.State.State)
}

func TestExtraDeviceService_UpdateDevice_WithMetadata(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	meta := json.RawMessage(`{"updated": true}`)
	device, err := svc.UpdateDevice(ctx2(), 1, "", "", meta)
	require.NoError(t, err)
	assert.NotNil(t, device.Metadata)
}

func TestIntegrationDeviceService_ListDevices_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	devices, total, err := svc.ListDevices(ctx2(), 1, 10, 999, nil, nil, "")
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, devices)
}

func TestIntegrationDeviceService_SetTags_InvalidKey(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.SetTags(ctx2(), "D001", map[string]string{"": "value"}, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid tag key")
}

func TestIntegrationDeviceService_SetTags_LongKey(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	longKey := make([]byte, 129)
	for i := range longKey {
		longKey[i] = 'a'
	}
	_, err := svc.SetTags(ctx2(), "D001", map[string]string{string(longKey): "value"}, true)
	assert.Error(t, err)
}

func TestIntegrationDeviceService_Tags_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.Tags(ctx2(), "D001")
	assert.Error(t, err)
}

func TestIntegrationDeviceService_RemoveTags_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	err := svc.RemoveTags(ctx2(), "D001", []string{"key"})
	assert.Error(t, err)
}

func TestIntegrationDeviceService_SetEnabled_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.SetEnabled(ctx2(), 999, true)
	assert.Error(t, err)
}

func TestIntegrationDeviceService_DeleteDevice_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	err := svc.DeleteDevice(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationDeviceService_RestoreDevice_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.RestoreDevice(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationDeviceService_Shadow_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.Shadow(ctx2(), "D001")
	assert.Error(t, err)
}

func TestIntegrationDeviceService_MutateShadow_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	desired := map[string]any{"key": "val"}
	_, err := svc.MutateShadow(ctx2(), "D001", 0, "app", &desired, nil, false)
	assert.Error(t, err)
}

func TestIntegrationDeviceService_ShadowHistory_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.ShadowHistory(ctx2(), "D001")
	assert.Error(t, err)
}

func TestIntegrationDeviceService_MQTT_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.MQTT(ctx2(), "D001")
	assert.Error(t, err)
}

func TestIntegrationDeviceService_SimulatePush_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.SimulatePush(ctx2(), "D001", `{"key":"val"}`, "test")
	assert.Error(t, err)
}

func TestIntegrationDeviceService_ListPushRecords_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	records, total, err := svc.ListPushRecords(ctx2(), "D001", 1, 10, "", "")
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, records)
}

func TestExtraDeviceService_PushRecord_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.PushRecord(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationDeviceService_ClearPushRecords_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	now := time.Now()
	count, err := svc.ClearPushRecords(ctx2(), "D001", &now)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestIntegrationDeviceService_Activate_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	_, err := svc.Activate(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationProductService_Get_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestProductService(db)

	_, err := svc.Get(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationProductService_GetByKey_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestProductService(db)

	_, err := svc.GetByKey(ctx2(), "NONEXIST")
	assert.Error(t, err)
}

func TestIntegrationProductService_Delete_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestProductService(db)

	err := svc.Delete(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationProductService_Restore_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestProductService(db)

	_, err := svc.Restore(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationRuleService_List_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestRuleService(db)

	rules, total, err := svc.List(ctx2(), 1, 10, "type1", "active", "search")
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, rules)
}

func TestIntegrationRuleService_Get_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestRuleService(db)

	_, err := svc.Get(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationRuleService_Delete_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestRuleService(db)

	err := svc.Delete(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationRuleService_ListExecutions_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestRuleService(db)

	executions, total, err := svc.ListExecutions(ctx2(), 999, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, executions)
}

func TestIntegrationAlertService_List_WithFilters(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestAlertService(db)

	alerts, total, err := svc.List(ctx2(), 1, 10, "critical", "high", "D001")
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, alerts)
}

func TestIntegrationOTAService_List_WithPagination(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestOTATotalService(db)

	packages, total, err := svc.List(ctx2(), 1, 5)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, packages)
}

func TestIntegrationOTAService_Batches_Empty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestOTATotalService(db)

	// Create a package first
	pkg := &model.OTAPackage{PackageName: "firmware-1", ProductKey: "P001", Version: "1.0.0"}
	err := svc.Create(ctx2(), pkg, "P001")
	require.NoError(t, err)

	batches, err := svc.Batches(ctx2(), pkg.PackageName)
	require.NoError(t, err)
	assert.Empty(t, batches)
}

func TestIntegrationOTAService_Deployments_Empty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestOTATotalService(db)

	// Create a package first
	pkg := &model.OTAPackage{PackageName: "firmware-1", ProductKey: "P001", Version: "1.0.0"}
	err := svc.Create(ctx2(), pkg, "P001")
	require.NoError(t, err)

	deployments, total, err := svc.Deployments(ctx2(), pkg.PackageName, 1, 10, "")
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, deployments)
}

func TestIntegrationOTAService_Delete_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestOTATotalService(db)

	err := svc.Delete(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationOTAService_Statistics_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestOTATotalService(db)

	_, err := svc.Statistics(ctx2(), "P001")
	assert.Error(t, err)
}

func TestIntegrationDeviceService_ListPushRecords_WithFilters(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	// Create push records
	_, err := svc.SimulatePush(ctx2(), "D001", `{"temp": 25}`, "test")
	require.NoError(t, err)

	// List with status filter
	records, total, err := svc.ListPushRecords(ctx2(), "D001", 1, 10, "", "")
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, records, 1)

	// List with operationType filter
	records, total, err = svc.ListPushRecords(ctx2(), "D001", 1, 10, "push", "")
	require.NoError(t, err)
	assert.Equal(t, int64(0), total) // No records with "push" operationType

	_ = records
}

func TestIntegrationDeviceService_BatchCreate_InvalidContent(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	// Invalid Excel content
	n, errs, err := svc.BatchCreate(ctx2(), []byte("not an excel file"))
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Empty(t, errs)
}

func TestIntegrationDeviceEventService_List_WithKeyword(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceEventService(db)

	events, total, err := svc.List(ctx2(), 1, 10, "test", "", "", nil, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	_ = events
}

func TestIntegrationDeviceEventService_List_WithDeviceKey(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceEventService(db)

	events, total, err := svc.List(ctx2(), 1, 10, "", "D001", "", nil, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	_ = events
}

func TestIntegrationDeviceEventService_List_WithEventType(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceEventService(db)

	events, total, err := svc.List(ctx2(), 1, 10, "", "", "temperature", nil, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	_ = events
}

func TestIntegrationDeviceEventService_Record_Valid(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceEventService(db)

	err := svc.Record(ctx2(), "P001", "D001", "temperature", 0, map[string]any{"temp": 25})
	require.NoError(t, err)

	events, total, err := svc.List(ctx2(), 1, 10, "", "D001", "temperature", nil, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, events, 1)
}

func TestIntegrationDeviceEventService_Record_WithTimestamp(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceEventService(db)

	err := svc.Record(ctx2(), "P001", "D001", "humidity", time.Now().UnixMilli(), map[string]any{"hum": 60})
	require.NoError(t, err)
}

func TestIntegrationDeviceEventService_Record_MissingFields(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceEventService(db)

	err := svc.Record(ctx2(), "", "D001", "temp", 0, nil)
	assert.Error(t, err)

	err = svc.Record(ctx2(), "P001", "", "temp", 0, nil)
	assert.Error(t, err)

	err = svc.Record(ctx2(), "P001", "D001", "", 0, nil)
	assert.Error(t, err)
}

func TestIntegrationDeviceEventService_Record_WrongProductKey(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceEventService(db)

	err := svc.Record(ctx2(), "WRONG", "D001", "temp", 0, nil)
	assert.Error(t, err)
}

func TestIntegrationDeviceEventService_Record_DeviceNotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceEventService(db)

	err := svc.Record(ctx2(), "P001", "NONEXIST", "temp", 0, nil)
	assert.Error(t, err)
}

func TestIntegrationRuleService_Evaluate(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	ruleRepo := repository.NewRuleRepository(&repository.IoTDB{DB: db})
	svc := service.NewRuleService(ruleRepo)

	// Create a rule first
	rule := &model.Rule{Name: "Test Rule", ProductID: 1, Type: "threshold", Status: "active"}
	err := svc.Create(ctx2(), rule)
	require.NoError(t, err)

	execution, err := svc.Evaluate(ctx2(), rule.ID)
	require.NoError(t, err)
	assert.NotNil(t, execution)
	assert.Equal(t, "success", execution.Status)
	assert.Equal(t, rule.ID, execution.RuleID)

	// Verify execution was saved
	executions, total, err := svc.ListExecutions(ctx2(), rule.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, executions, 1)

	// Verify rule stats updated
	updatedRule, err := svc.Get(ctx2(), rule.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), updatedRule.ExecutionCount)
	assert.Equal(t, int64(1), updatedRule.SuccessCount)
}

func TestIntegrationRuleService_Evaluate_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestRuleService(db)

	_, err := svc.Evaluate(ctx2(), 999)
	assert.Error(t, err)
}

func TestIntegrationAlertService_SetStatus_Acknowledged(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestAlertService(db)

	rule := &model.Rule{Name: "Test Rule", ProductID: 1, Type: "threshold", Status: "active"}
	err := db.Create(rule).Error
	require.NoError(t, err)

	alert := &model.Alert{RuleID: rule.ID, RuleName: rule.Name, Status: "active", Severity: "critical"}
	err = db.Create(alert).Error
	require.NoError(t, err)

	result, err := svc.SetStatus(ctx2(), alert.ID, "acknowledged")
	require.NoError(t, err)
	assert.Equal(t, "acknowledged", result.Status)
	assert.NotNil(t, result.AckAt)
	assert.Nil(t, result.ResolvedAt)
}

func TestIntegrationAlertService_SetStatus_Resolved(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestAlertService(db)

	rule := &model.Rule{Name: "Test Rule", ProductID: 1, Type: "threshold", Status: "active"}
	err := db.Create(rule).Error
	require.NoError(t, err)

	alert := &model.Alert{RuleID: rule.ID, RuleName: rule.Name, Status: "active", Severity: "critical"}
	err = db.Create(alert).Error
	require.NoError(t, err)

	result, err := svc.SetStatus(ctx2(), alert.ID, "resolved")
	require.NoError(t, err)
	assert.Equal(t, "resolved", result.Status)
	assert.NotNil(t, result.ResolvedAt)
	assert.Nil(t, result.AckAt)
}

func TestIntegrationAlertService_SetStatus_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestAlertService(db)

	_, err := svc.SetStatus(ctx2(), 999, "resolved")
	assert.Error(t, err)
}

func TestIntegrationOTAService_Statistics_WithData(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestOTATotalService(db)

	// Create a package
	pkg := &model.OTAPackage{PackageName: "fw-1.0", Version: "1.0.0"}
	err := svc.Create(ctx2(), pkg, "P001")
	require.NoError(t, err)

	stats, err := svc.Statistics(ctx2(), pkg.PackageName)
	require.NoError(t, err)
	assert.Equal(t, "1", stats.PackageID)
	assert.Equal(t, int64(0), stats.TotalTargetDevices)
}

func TestIntegrationOTAService_Update_Valid(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestOTATotalService(db)

	pkg := &model.OTAPackage{PackageName: "fw-2.0", Version: "2.0.0"}
	err := svc.Create(ctx2(), pkg, "P001")
	require.NoError(t, err)

	pkg.Version = "2.1.0"
	err = svc.Update(ctx2(), pkg)
	require.NoError(t, err)

	fetched, err := svc.Get(ctx2(), pkg.ID)
	require.NoError(t, err)
	assert.Equal(t, "2.1.0", fetched.Version)
}

func TestIntegrationOTAService_Update_ProductNotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestOTATotalService(db)

	pkg := &model.OTAPackage{ID: 1, ProductID: 999}
	err := svc.Update(ctx2(), pkg)
	assert.Error(t, err)
}

func TestIntegrationProductService_Create_InvalidMetadata(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestProductService(db)

	p := &model.Product{Name: "Bad Meta", Metadata: json.RawMessage(`"invalid-string-not-json"`)}
	_, err := svc.Create(ctx2(), p)
	assert.Error(t, err)
}

func TestIntegrationProductService_Create_Defaults(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestProductService(db)

	p := &model.Product{Name: "Defaults"}
	result, err := svc.Create(ctx2(), p)
	require.NoError(t, err)
	assert.NotEmpty(t, result.ProductKey)
	assert.Equal(t, "active", result.Status)
	assert.Equal(t, "direct", result.NodeType)
}

func TestIntegrationProductService_Create_WithKey(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestProductService(db)

	p := &model.Product{Name: "With Key", ProductKey: "CUSTOM"}
	result, err := svc.Create(ctx2(), p)
	require.NoError(t, err)
	assert.Equal(t, "CUSTOM", result.ProductKey)
}

func TestIntegrationProductService_Save_InvalidMetadata(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestProductService(db)

	p := &model.Product{ID: 1, Name: "Updated", ProductKey: "P001", Metadata: json.RawMessage(`"bad"`)}
	err := svc.Save(ctx2(), p)
	assert.Error(t, err)
}

func TestIntegrationProductService_Delete_HasDevices(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestProductService(db)

	err := svc.Delete(ctx2(), 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "has devices")
}

func TestIntegrationDeviceService_BatchCreate_ValidExcel(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	var buf bytes.Buffer
	f := excelize.NewFile()
	f.SetSheetRow("Sheet1", "A1", &[]string{"ProductKey", "DeviceName"})
	f.SetSheetRow("Sheet1", "A2", &[]string{"P001", "Batch Device 1"})
	f.SetSheetRow("Sheet1", "A3", &[]string{"P001", "Batch Device 2"})
	f.Write(&buf)

	n, errs, err := svc.BatchCreate(ctx2(), buf.Bytes())
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Empty(t, errs)
}

func TestIntegrationDeviceService_BatchCreate_NoRows(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	var buf bytes.Buffer
	f := excelize.NewFile()
	f.SetSheetRow("Sheet1", "A1", &[]string{"ProductKey", "DeviceName"})
	f.Write(&buf)

	n, errs, err := svc.BatchCreate(ctx2(), buf.Bytes())
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Empty(t, errs)
}

func TestIntegrationDeviceService_BatchCreate_MissingColumns(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	var buf bytes.Buffer
	f := excelize.NewFile()
	f.SetSheetRow("Sheet1", "A1", &[]string{"Col1", "Col2"})
	f.SetSheetRow("Sheet1", "A2", &[]string{"val1", "val2"})
	f.Write(&buf)

	n, errs, err := svc.BatchCreate(ctx2(), buf.Bytes())
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Empty(t, errs)
}

func TestIntegrationDeviceService_BatchCreate_WrongProduct(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	var buf bytes.Buffer
	f := excelize.NewFile()
	f.SetSheetRow("Sheet1", "A1", &[]string{"ProductKey", "DeviceName"})
	f.SetSheetRow("Sheet1", "A2", &[]string{"WRONG", "Device 1"})
	f.Write(&buf)

	n, errs, err := svc.BatchCreate(ctx2(), buf.Bytes())
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Len(t, errs, 1)
}

func TestIntegrationDeviceService_BatchCreate_EmptyRow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestDeviceService(db)

	var buf bytes.Buffer
	f := excelize.NewFile()
	f.SetSheetRow("Sheet1", "A1", &[]string{"ProductKey", "DeviceName"})
	f.SetSheetRow("Sheet1", "A2", &[]string{"", ""})
	f.SetSheetRow("Sheet1", "A3", &[]string{"P001", "Good Device"})
	f.Write(&buf)

	n, errs, err := svc.BatchCreate(ctx2(), buf.Bytes())
	require.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Empty(t, errs)
}

func TestExtraDeviceService_MockKafka_NoBrokers(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestDeviceService(db)

	err := svc.MockKafka(ctx2(), nil, "topic", "data")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}

func TestIntegrationProductTSLService_Upsert_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestProductTSLService(db)

	err := svc.Upsert(ctx2(), "P001", `{"properties":[]}`)
	require.NoError(t, err)

	result, err := svc.Get(ctx2(), "P001")
	require.NoError(t, err)
	assert.Equal(t, `{"properties":[]}`, result.TSL)
}

func TestIntegrationProductTSLService_Upsert_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestProductTSLService(db)

	err := svc.Upsert(ctx2(), "P001", `{"properties":[]}`)
	require.NoError(t, err)

	err = svc.Upsert(ctx2(), "P001", `{"properties":[{"name":"temp"}]}`)
	require.NoError(t, err)

	result, err := svc.Get(ctx2(), "P001")
	require.NoError(t, err)
	assert.Equal(t, `{"properties":[{"name":"temp"}]}`, result.TSL)
}

func TestIntegrationProductTSLService_Get_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestProductTSLService(db)

	_, err := svc.Get(ctx2(), "NONEXIST")
	assert.Error(t, err)
}

func TestIntegrationProductTSLService_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.SeedTestData(t, db)
	svc := testutil.NewTestProductTSLService(db)

	err := svc.Upsert(ctx2(), "P001", `{"properties":[]}`)
	require.NoError(t, err)

	err = svc.Delete(ctx2(), "P001")
	require.NoError(t, err)

	_, err = svc.Get(ctx2(), "P001")
	assert.Error(t, err)
}

func TestIntegrationProductTSLService_Delete_ProductNotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestProductTSLService(db)

	err := svc.Delete(ctx2(), "NONEXIST")
	assert.Error(t, err)
}

func TestIntegrationProductTSLService_Upsert_ProductNotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := testutil.NewTestProductTSLService(db)

	err := svc.Upsert(ctx2(), "NONEXIST", `{"properties":[]}`)
	assert.Error(t, err)
}
