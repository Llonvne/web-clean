package database

import (
	"gorm.io/gorm"
)

var models = make([]any, 0)

func RegisterSchema(model any) {
	models = append(models, model)
}

func AutoMigrateRegisteredSchema(database Database) error {
	return database.Transaction(func(tx *gorm.DB) error {
		return tx.AutoMigrate(models...)
	})
}
