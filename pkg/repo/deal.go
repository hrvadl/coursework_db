package repo

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
)

type Deal interface {
	GetByID(id int) (*models.Deal, error)
	Get() ([]models.Deal, error)
}

func NewDeal(db *gorm.DB) Deal {
	return &deal{db}
}

type deal struct {
	db *gorm.DB
}

func (d *deal) GetByID(id int) (*models.Deal, error) {
	var deal models.Deal

	if err := d.db.Find(&deal, id).Error; err != nil {
		return nil, err
	}

	return &deal, nil
}

func (d *deal) Get() ([]models.Deal, error) {
	var deals []models.Deal

	if err := d.db.Find(&deals).Error; err != nil {
		return nil, err
	}

	return deals, nil
}
