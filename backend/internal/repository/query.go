package repository

import (
	"aiot-backend/internal/dal/query"
	"gorm.io/gorm"
)

func useIoTQuery(db *IoTDB) *query.Query {
	return query.Use(db.DB)
}

func useQuery(db *gorm.DB) *query.Query {
	return query.Use(db)
}
