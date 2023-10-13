package services

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type Emitent interface {
	GetByID(id int) (*models.User, error)
	Get() ([]models.User, error)
}

type emitent struct {
	repo   repo.Emitent
	crypto Cryptor
}

func NewEmitent(repo repo.Emitent, crypto Cryptor) Emitent {
	return &emitent{
		repo:   repo,
		crypto: crypto,
	}
}

func (s *emitent) Get() ([]models.User, error) {
	return s.repo.Get()
}

func (s *emitent) GetByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}
