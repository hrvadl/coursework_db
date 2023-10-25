package repo

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Inventory interface {
	GetByID(id int) (*models.InventoryItem, error)
	GetByUserID(userID int) ([]models.InventoryItem, error)
	GetUserInventoryBySecurityID(userID int, securityID int) (*models.InventoryItem, error)
	Patch(item *models.InventoryItem) (*models.InventoryItem, error)
	Save(item *models.InventoryItem) error
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

func (i *inventory) GetByUserID(userID int) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	err := i.db.Where(&models.InventoryItem{OwnerID: uint(userID)}).Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (i *inventory) GetUserInventoryBySecurityID(userID int, securityID int) (*models.InventoryItem, error) {
	var inventory models.InventoryItem
	res := i.db.Model(&models.InventoryItem{}).
		Preload(clause.Associations).
		Where(&models.InventoryItem{SecurityID: uint(securityID), OwnerID: uint(userID)}).
		First(&inventory)

	if res.Error != nil {
		return nil, res.Error
	}

	return &inventory, nil
}

func (i *inventory) Save(item *models.InventoryItem) error {
	return i.db.Save(item).Error
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
