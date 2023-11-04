package db

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

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

	result := db.First(&models.Security{})
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return db, nil
	}

	if err := seed(db); err != nil {
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

func seed(db *gorm.DB) error {
	content, err := os.ReadFile(path.Join("../", "seed.sql"))

	if err != nil {
		log.Printf("Error reading seed.sql: %v \n", err)
		return err
	}

	sql := string(content)
	res := db.Exec(sql)

	if err := res.Error; err != nil {
		log.Printf("Error seeding db: %v \n", err)
		return err
	}

	if res.RowsAffected == 0 {
		log.Println("Zero rows affected after seeding db")
		return errors.New("seed failed: no rows affected")
	}

	return nil
}
