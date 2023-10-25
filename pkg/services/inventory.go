package services

import (
	"errors"

	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type Inventory interface {
	GetByID(id int) (*models.InventoryItem, error)
	GetUserInventory(userID int) ([]models.InventoryItem, error)
	GetUserInventoryBySecurityID(userID int, securityID int) (*models.InventoryItem, error)
	CreateOrUpdate(item *models.InventoryItem) (*models.InventoryItem, error)
	Patch(item *models.InventoryItem) (*models.InventoryItem, error)
	Delete(id int) error
}

func NewInventory(repo repo.Inventory, ur repo.User) Inventory {
	return &inventory{
		repo: repo,
		user: ur,
	}
}

type inventory struct {
	repo repo.Inventory
	user repo.User
}

func (i *inventory) GetByID(id int) (*models.InventoryItem, error) {
	return i.repo.GetByID(id)
}

func (i *inventory) GetUserInventory(userID int) ([]models.InventoryItem, error) {
	return i.repo.GetByUserID(userID)
}

func (i *inventory) GetUserInventoryBySecurityID(userID int, securityID int) (*models.InventoryItem, error) {
	return i.repo.GetUserInventoryBySecurityID(userID, securityID)
}

func (i *inventory) CreateOrUpdate(item *models.InventoryItem) (*models.InventoryItem, error) {
	if err := i.repo.Save(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (i *inventory) Patch(item *models.InventoryItem) (*models.InventoryItem, error) {
	if item.Amount <= 0 {
		return nil, errors.New("amount of securities cannot be less than 0")
	}

	if item.ID == 0 {
		return nil, errors.New("item id is not specified")
	}

	return i.repo.Patch(item)
}

func (i *inventory) Delete(id int) error {
	return i.repo.Delete(id)
}
