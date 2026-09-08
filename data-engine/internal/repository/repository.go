package repository

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Repository is the base repository holding the database connection.
type Repository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRepository(logger *zap.Logger, db *gorm.DB) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// DB returns the underlying gorm.DB instance.
func (r *Repository) DB() *gorm.DB {
	return r.db
}

// NewDB initializes a GORM DB connection from configuration.
func NewDB(config *viper.Viper, logger *zap.Logger) (*gorm.DB, error) {
	driver := config.GetString("data.db.aiot.driver")
	dsn := config.GetString("data.db.aiot.dsn")
	if driver == "" || dsn == "" {
		return nil, fmt.Errorf("data.db.aiot.driver and data.db.aiot.dsn are required")
	}

	var dialect gorm.Dialector
	switch driver {
	case "mysql":
		dialect = mysql.Open(dsn)
	case "postgres", "postgresql":
		dialect = postgres.Open(dsn)
	case "sqlite":
		dialect = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported ota database driver %q", driver)
	}

	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return db, nil
}
