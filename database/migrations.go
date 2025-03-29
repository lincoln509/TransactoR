package database

import (
	"fmt"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB, models []interface{}) error {
	if db == nil {
		return fmt.Errorf("database non initialisée")
	}

	if db.Config.Dialector.Name() == "sqlite" {
		db.Exec("PRAGMA foreign_keys = OFF")
		defer db.Exec("PRAGMA foreign_keys = ON")
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("échec migration %T: %w", model, err)
		}
	}

	return nil
}
