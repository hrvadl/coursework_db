package db

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Deal{},
		&models.InventoryItem{},
		&models.Security{},
		&models.Transaction{},
		&models.User{},
		&models.Session{},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Must(db *gorm.DB, err error) *gorm.DB {
	if err != nil {
		panic(err)
	}

	return db
}
