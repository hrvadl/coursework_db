package repo

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
)

type Deal interface {
	GetByID(id int) (*models.Deal, error)
	Get() ([]models.Deal, error)
	Patch(deal *models.Deal) (*models.Deal, error)
	Delete(id int) error
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

func (d *deal) Patch(deal *models.Deal) (*models.Deal, error) {
	if err := d.db.Model(deal).Updates(deal).Error; err != nil {
		return nil, err
	}

	return deal, nil
}

func (d *deal) Delete(id int) error {
	_, err := d.Patch(&models.Deal{ID: uint(id), Active: false})
	return err
}
