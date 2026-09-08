package repository

import (
	"data-engine/internal/dal/query"

	"gorm.io/gorm"
)

func useQuery(db *gorm.DB) *query.Query {
	return query.Use(db)
}
