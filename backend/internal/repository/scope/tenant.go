package scope

import (
	"aiot-backend/internal/tenant"

	"gorm.io/gorm"
)

// Tenant filters the current query by the tenant id resolved from the DB
// statement's context. The column is resolved as "<table>.tenant_id" from the
// statement's model/table; when undetermined it falls back to "tenant_id"
// (safe for single-table queries). Use it directly with GORM's Scopes:
//
//	db.Scopes(scope.Tenant).Find(&products)
func Tenant(db *gorm.DB) *gorm.DB {
	table := db.Statement.Table
	if table == "" && db.Statement.Schema != nil {
		table = db.Statement.Schema.Table
	}
	column := "tenant_id"
	if table != "" {
		column = table + ".tenant_id"
	}
	return db.Where(column+" = ?", tenant.GetTenantID(db.Statement.Context))
}
