package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Driver   string
	DSN      string
	LogLevel logger.LogLevel
}

var dbInstance *gorm.DB

func InitDB(cfg Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "postgres":
		dialector = postgres.Open(cfg.DSN)
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("driver non supporté")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	})

	if err != nil {
		return nil, fmt.Errorf("échec de connexion DB: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	dbInstance = db
	return db, nil
}

func GetDB() *gorm.DB {
	return dbInstance
}

// InitDBWithMigrations - Initialisation DB avec migrations automatiques
func InitDBWithMigrations(cfg Config, models []interface{}) (*gorm.DB, error) {
	db, err := InitDB(cfg)
	if err != nil {
		return nil, err
	}

	if err := Migrate(db, models, MigrationConfig{}); err != nil {
		return nil, fmt.Errorf("échec des migrations: %v", err)
	}

	return db, nil
}
