// Package database provides SQLite database initialization and migration.
package database

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitSQLite opens (or creates) the SQLite database file configured in the YAML
// and returns a *gorm.DB ready for use.
func InitSQLite(cfg config.TypeMyPortfolio) *gorm.DB {
	dsn := cfg.Database.DSN

	// Ensure the parent directory exists.
	dir := filepath.Dir(dsn)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create database directory %s: %v", dir, err)
	}

	logLevel := logger.Warn
	switch cfg.Database.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "info":
		logLevel = logger.Info
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	// Enable WAL mode for better concurrency with SQLite.
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	log.Printf("Database connected: %s", dsn)
	return db
}

// AutoMigrate runs GORM auto-migration for every model.
func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&model.Admin{},
		&model.Owner{},
		&model.Project{},
		&model.Experience{},
		&model.Skill{},
		&model.SocialLink{},
		&model.UploadedFile{},
		&model.OAuthUser{},
		&model.Comment{},
		&model.ContactMessage{},
		&model.TechStack{},
		&model.Post{},
		&model.UpcomingItem{},
		&model.Certificate{},
	); err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}
	log.Println("Database migration complete")
}
