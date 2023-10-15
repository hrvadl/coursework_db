package services

import (
	"errors"

	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type Deal interface {
	GetByID(id int) (*models.Deal, error)
	Get() ([]models.Deal, error)
	Create(d *models.Deal) (*models.Deal, error)
}

func NewDeal(repo repo.Deal) Deal {
	return &deal{repo}
}

type deal struct {
	repo repo.Deal
}

func (d *deal) GetByID(id int) (*models.Deal, error) {
	return d.repo.GetByID(id)
}

func (d *deal) Get() ([]models.Deal, error) {
	return d.repo.Get()
}

func (d *deal) Create(deal *models.Deal) (*models.Deal, error) {
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
