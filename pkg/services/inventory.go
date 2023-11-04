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
	Withdraw(ownerID, securityID, amount uint) (*models.InventoryItem, error)
	Add(ownerID, securityID, amount uint) (*models.InventoryItem, error)
	Patch(item *models.InventoryItem) (*models.InventoryItem, error)
	Delete(id int) error
}

func NewInventory(repo repo.Inventory, ur repo.User, dr repo.Deal) Inventory {
	return &inventory{
		repo: repo,
		user: ur,
		deal: dr,
	}
}

type inventory struct {
	repo repo.Inventory
	user repo.User
	deal repo.Deal
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

func (i *inventory) Withdraw(ownerID, securityID, amount uint) (*models.InventoryItem, error) {
	if amount < 1 {
		return nil, errors.New("amount must be greater than zero")
	}

	inventory, err := i.GetUserInventoryBySecurityID(int(ownerID), int(securityID))

	if err != nil {
		return nil, err
	}

	if inventory == nil {
		return nil, errors.New("you have nothing to withdraw")
	}

	if int(inventory.Amount)-int(amount) < 0 {
		return nil, errors.New("total amount must be greater than zero")
	}

	inventory.Amount -= amount

	if err := i.repo.Save(inventory); err != nil {
		return nil, err
	}

	hasDeal, err := i.deal.GetByUserIDAndSecurityID(int(ownerID), int(securityID))

	if hasDeal.Amount <= inventory.Amount {
		return inventory, nil
	}

	if err != nil {
		return inventory, nil
	}

	if amount >= hasDeal.Amount {
		err := i.deal.Delete(int(hasDeal.ID))
		return inventory, err
	}

	hasDeal.Amount -= amount
	if _, err := i.deal.Patch(hasDeal); err != nil {
		return inventory, err
	}

	return inventory, nil
}

func (i *inventory) Add(ownerID, securityID, amount uint) (*models.InventoryItem, error) {
	if amount < 1 {
		return nil, errors.New("amount must be greater than zero")
	}

	inventory, _ := i.GetUserInventoryBySecurityID(int(ownerID), int(securityID))

	if inventory == nil {
		inventory = &models.InventoryItem{}
	}

	inventory.OwnerID = ownerID
	inventory.Amount += amount
	inventory.SecurityID = securityID

	if err := i.repo.Save(inventory); err != nil {
		return nil, err
	}

	return inventory, nil
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
