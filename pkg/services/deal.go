package services

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type Deal interface {
	GetByID(id int) (*models.Deal, error)
	Get() ([]models.Deal, error)
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
