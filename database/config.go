package database

import (
	"fmt"
	"safelyyou_coding_challenge_go/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// InitDatabase connects to SQLite and configures the database
func InitDatabase(dbUrl string) (*DatabaseOperations, error) {
	db, err := connectSQLite(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Exec("PRAGMA journal_mode = WAL;").Error; err != nil {
		return nil, err
	}

	if err := db.Exec("PRAGMA synchronous = NORMAL;").Error; err != nil {
		return nil, err
	}

	if err := configureConnectionPool(db); err != nil {
		return nil, fmt.Errorf("failed to configure connection pool: %w", err)
	}

	// Auto-migrate database tables
	if err := migrateDatabase(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &DatabaseOperations{DB: db}, nil
}

// connectSQLite creates a GORM SQLite connection
func connectSQLite(dbUrl string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dbUrl,
	}, &gorm.Config{})
}

// configureConnectionPool sets up SQLite connection pool settings
func configureConnectionPool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Optimize SQLite connection pool for performance
	sqlDB.SetMaxOpenConns(1)            // Allow more concurrent connections
	sqlDB.SetMaxIdleConns(1)            // Keep connections alive
	sqlDB.SetConnMaxLifetime(time.Hour) // Connection lifetime

	return nil
}

// migrateDatabase creates all necessary database tables
func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(models.AllModels()...)
}
