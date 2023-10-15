package services

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type Security interface {
	Get() ([]models.Security, error)
}

func NewSecurity(sr repo.Security) Security {
	return &security{sr: sr}
}

type security struct {
	sr repo.Security
}

func (s *security) Get() ([]models.Security, error) {
	return s.sr.Get()
}
