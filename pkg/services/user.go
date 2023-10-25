package services

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type User interface {
	GetByID(id int) (*models.User, error)
	Get() ([]models.User, error)
	Patch(user *models.User) (*models.User, error)
}

type stock struct {
	repo   repo.User
	crypto Cryptor
}

func NewStock(repo repo.User, crypto Cryptor) User {
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

func (s *stock) Patch(user *models.User) (*models.User, error) {
	return s.repo.Patch(user)
}
