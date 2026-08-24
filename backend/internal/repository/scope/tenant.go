package scope

import (
	"aiot-backend/internal/tenant"

	"gorm.io/gorm"
)

// Tenant filters the current query by the tenant id resolved from the DB
// statement's context. The column is resolved as "<table>.organization_id" from the
// statement's model/table; when undetermined it falls back to "organization_id"
// (safe for single-table queries). Use it directly with GORM's Scopes:
//
//	db.Scopes(scope.Tenant).Find(&products)
func Tenant(db *gorm.DB) *gorm.DB {
	table := db.Statement.Table
	if table == "" && db.Statement.Schema != nil {
		table = db.Statement.Schema.Table
	}
	column := "organization_id"
	if table != "" {
		column = table + ".organization_id"
	}
	return db.Where(column+" = ?", tenant.GetOrganizationID(db.Statement.Context))
}
