package repo

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
)

type Deal interface {
	GetAllByID(id int) (*models.Deal, error)
	GetByID(id int) (*models.Deal, error)
	Get() ([]models.Deal, error)
	GetByUserIDAndSecurityID(userID, securityID int) (*models.Deal, error)
	Create(deal *models.Deal) (*models.Deal, error)
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
	err := d.db.Model(&models.Deal{}).
		Where(&models.Deal{Active: true}).
		Preload("Owner").
		Preload("Security").
		First(&deal, id).Error

	if err != nil {
		return nil, err
	}

	return &deal, nil
}

func (d *deal) GetAllByID(id int) (*models.Deal, error) {
	var deal models.Deal
	err := d.db.Model(&models.Deal{}).
		Preload("Owner").
		Preload("Security").
		First(&deal, id).Error

	if err != nil {
		return nil, err
	}

	return &deal, nil
}

func (d *deal) Get() ([]models.Deal, error) {
	var deals []models.Deal
	res := d.db.Model(&models.Deal{}).
		Where(&models.Deal{Active: true}).
		Preload("Security").
		Find(&deals)

	if err := res.Error; err != nil {
		return nil, err
	}

	return deals, nil
}

func (d *deal) GetByUserIDAndSecurityID(userID, securityID int) (*models.Deal, error) {
	var deal models.Deal
	res := d.db.Model(&models.Deal{}).
		Where(&models.Deal{
			Active:     true,
			SecurityID: uint(securityID),
			OwnerID:    uint(userID),
		}).
		First(&deal)

	return &deal, res.Error
}

func (d *deal) Patch(deal *models.Deal) (*models.Deal, error) {
	if err := d.db.Model(deal).Updates(deal).Error; err != nil {
		return nil, err
	}

	return deal, nil
}

func (d *deal) Create(deal *models.Deal) (*models.Deal, error) {
	if err := d.db.Create(deal).Error; err != nil {
		return nil, err
	}

	return deal, nil
}

func (d *deal) Delete(id int) error {
	res := d.db.Model(&models.Deal{}).
		Where("id = ?", id).
		Update("active", false)

	return res.Error
}
