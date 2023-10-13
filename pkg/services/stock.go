package services

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type Stock interface {
	GetByID(id int) (*models.User, error)
	Get() ([]models.User, error)
}

type stock struct {
	repo   repo.Stock
	crypto Cryptor
}

func NewStock(repo repo.Stock, crypto Cryptor) Stock {
	return &stock{
		repo:   repo,
		crypto: crypto,
	}
}

func (s *stock) Get() ([]models.User, error) {
	return s.repo.Get()
}

func (s *stock) GetByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}
