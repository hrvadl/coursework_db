package services

import (
	"errors"

	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type Deal interface {
	MakeTransaction(authorID, dealID, amount int) error
	GetByID(id int) (*models.Deal, error)
	Get() ([]models.Deal, error)
	Create(d *models.Deal) (*models.Deal, error)
	Patch(d *models.Deal) (*models.Deal, error)
	Delete(d int) error
}

func NewDeal(
	repo repo.Deal,
	irepo repo.Inventory,
	urepo repo.User,
) Deal {
	return &deal{
		repo:  repo,
		irepo: irepo,
		urepo: urepo,
	}
}

type deal struct {
	repo  repo.Deal
	irepo repo.Inventory
	urepo repo.User
}

func (d *deal) GetByID(id int) (*models.Deal, error) {
	return d.repo.GetByID(id)
}

func (d *deal) Get() ([]models.Deal, error) {
	return d.repo.Get()
}

func (d *deal) Create(deal *models.Deal) (*models.Deal, error) {
	hasSecurity, _ := d.irepo.GetUserInventoryBySecurityID(int(deal.OwnerID), int(deal.SecurityID))

	if deal.Sell && hasSecurity == nil {
		return nil, errors.New("cannot sell deal without security")
	}

	if !deal.Active {
		return nil, errors.New("cannot create inactive deal")
	}

	if deal.Amount <= 0 {
		return nil, errors.New("cannot create a deal with a negative or zero amount")
	}

	if deal.SecurityID == 0 {
		return nil, errors.New("securityID cannot be empty")
	}

	if deal.OwnerID == 0 {
		return nil, errors.New("ownerID cannot be empty")
	}

	if deal.Price <= 0 {
		return nil, errors.New("cannot create a deal with a negative or zero amount")
	}

	return d.repo.Create(deal)
}

func (d *deal) Patch(deal *models.Deal) (*models.Deal, error) {
	if deal.ID == 0 {
		return nil, errors.New("deal cannot be empty")
	}

	return d.repo.Patch(deal)
}

func (d *deal) Delete(id int) error {
	err := d.repo.Delete(id)
	return err
}

func (d *deal) MakeTransaction(authorID, dealID, amount int) error {
	deal, err := d.GetByID(dealID)

	if err != nil {
		return err
	}

	if !deal.Active {
		return errors.New("deal is not active")
	}

	if authorID == int(deal.OwnerID) {
		return errors.New("cannot make transaction with yourself")
	}

	if amount < 0 {
		return errors.New("amount must be greater than zero")
	}

	switch deal.Sell {
	}

	return nil
}

func (d *deal) buy()
