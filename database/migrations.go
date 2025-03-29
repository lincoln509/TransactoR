package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// MigrationLogger - Logger personnalisé pour les migrations
type MigrationLogger struct {
	*log.Logger
}

func (ml *MigrationLogger) Printf(format string, v ...interface{}) {
	ml.Output(2, fmt.Sprintf("[Migration] "+format, v...))
}

// MigrationConfig - Configuration des migrations
type MigrationConfig struct {
	DropIndexes    bool
	DisableLogging bool
}

// Migrate - Exécute les migrations pour tous les modèles
func Migrate(db *gorm.DB, models []interface{}, config MigrationConfig) error {
	// Configuration du logger
	if !config.DisableLogging {
		db.Logger = &MigrationLogger{
			Logger: log.New(log.Writer(), "", log.LstdFlags|log.Lshortfile),
		}
	}

	// Désactiver les contraintes FK pour SQLite
	if db.Config.Dialector.Name() == "sqlite" {
		if err := db.Exec("PRAGMA foreign_keys = OFF").Error; err != nil {
			return fmt.Errorf("échec désactivation FK: %v", err)
		}
	}

	// Migration des modèles
	for _, model := range models {
		if config.DropIndexes {
			if err := dropIndexes(db, model); err != nil {
				return err
			}
		}

		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("échec migration %T: %v", model, err)
		}
	}

	// Réactiver les FK pour SQLite
	if db.Config.Dialector.Name() == "sqlite" {
		if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			return fmt.Errorf("échec réactivation FK: %v", err)
		}
	}

	return nil
}

// dropIndexes - Supprime les index existants avant migration
func dropIndexes(db *gorm.DB, model interface{}) error {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return err
	}

	rows, err := db.Migrator().GetIndexes(stmt.Table)
	if err != nil {
		return err
	}

	for _, index := range rows {
		if index.Name != "PRIMARY" { // Ne pas supprimer la clé primaire
			if err := db.Migrator().DropIndex(model, index.Name); err != nil {
				return fmt.Errorf("échec suppression index %s: %v", index.Name, err)
			}
		}
	}

	return nil
}

// MigrationExample - Exemple d'utilisation
// func MigrationExample() {
//     db := GetDB()

//     models := []interface{}{
//         &User{},
//         &Product{},
//         &Order{},
//     }

//     config := MigrationConfig{
//         DropIndexes:    true,
//         DisableLogging: false,
//     }

//     if err := Migrate(db, models, config); err != nil {
//         log.Fatalf("Échec des migrations: %v", err)
//     }
// }
