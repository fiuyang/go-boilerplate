package migrations

import (
	"gin-boilerplate/model"

	"gorm.io/gorm"
)

// AutoMigrateIfNotExists migrates the model if the table doesn't exist
func AutoMigrateIfNotExists(db *gorm.DB, modelInstance interface{}) {
	if !db.Migrator().HasTable(modelInstance) {
		db.AutoMigrate(modelInstance)
	} 
}

func PerformMigrations(db *gorm.DB) {
	AutoMigrateIfNotExists(db, &model.Users{})
	AutoMigrateIfNotExists(db, &model.Tags{})
	AutoMigrateIfNotExists(db, &model.PasswordResets{})
}