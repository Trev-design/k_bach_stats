package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitialMigration(schemas ...any) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("messages.db"))
	if err != nil {
		return nil, err
	}

	for _, schema := range schemas {
		if err := migrateSchemas(db, schema); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func migrateSchemas(db *gorm.DB, schema any) error {
	if err := db.AutoMigrate(schema); err != nil {
		return err
	}

	return nil
}
