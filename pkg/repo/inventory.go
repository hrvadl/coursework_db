package repo

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
)

type Inventory interface {
	GetByID(id int) (*models.InventoryItem, error)
	Get(userID int) ([]models.InventoryItem, error)
	Patch(item *models.InventoryItem) (*models.InventoryItem, error)
	Delete(id int) error
}

func NewInventory(db *gorm.DB) Inventory {
	return &inventory{db}
}

type inventory struct {
	db *gorm.DB
}

func (i *inventory) GetByID(id int) (*models.InventoryItem, error) {
	var item models.InventoryItem

	if err := i.db.Find(&item, id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (i *inventory) Get(userID int) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	err := i.db.Where(&models.InventoryItem{OwnerID: uint(userID)}).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (i *inventory) Patch(item *models.InventoryItem) (*models.InventoryItem, error) {
	if err := i.db.Model(item).Updates(item).Error; err != nil {
		return nil, err
	}

	return item, nil
}

func (i *inventory) Delete(id int) error {
	res := i.db.Delete(&models.InventoryItem{}, id)
	return res.Error
}
