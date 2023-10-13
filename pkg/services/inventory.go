package services

import (
	"errors"

	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type Inventory interface {
	GetByID(id int) (*models.InventoryItem, error)
	GetUserInventory(userID int) ([]models.InventoryItem, error)
	Patch(item *models.InventoryItem) (*models.InventoryItem, error)
	Delete(id int) error
}

func NewInventory(repo repo.Inventory, sr repo.Stock, er repo.Emitent) Inventory {
	return &inventory{
		repo:        repo,
		stockRepo:   sr,
		emitentRepo: er,
	}
}

type inventory struct {
	repo        repo.Inventory
	stockRepo   repo.Stock
	emitentRepo repo.Emitent
}

func (i *inventory) GetByID(id int) (*models.InventoryItem, error) {
	return i.repo.GetByID(id)
}

func (i *inventory) GetUserInventory(userID int) ([]models.InventoryItem, error) {
	return i.repo.Get(userID)
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
