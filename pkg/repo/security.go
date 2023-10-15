package repo

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
)

type Security interface {
	Get() ([]models.Security, error)
}

func NewSecurity(db *gorm.DB) Security {
	return &security{db: db}
}

type security struct {
	db *gorm.DB
}

func (s *security) Get() ([]models.Security, error) {
	var sec []models.Security

	if err := s.db.Find(&sec).Error; err != nil {
		return nil, err
	}

	return sec, nil
}
