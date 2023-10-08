package repo

import (
	"errors"
	"fmt"
	"time"

	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
)

type Session interface {
	Create(id uint) (*models.Session, error)
	GetByUserID(id uint) (*models.Session, error)
}

func NewSession(db *gorm.DB) Session {
	return &session{
		db: db,
	}
}

type session struct {
	db *gorm.DB
}

func (s *session) Create(id uint) (*models.Session, error) {
	session := &models.Session{
		UserID:     id,
		ValidUntil: time.Now().Add(time.Hour * 24),
	}

	_ = s.DeleteByUserID(id)
	res := s.db.Create(session)

	if res.Error != nil {
		return nil, res.Error
	}

	return session, nil
}

func (s *session) GetByUserID(id uint) (*models.Session, error) {
	var session models.Session
	err := s.db.Where(&models.Session{UserID: id}, id).First(&session).Error

	if err != nil {
		return nil, fmt.Errorf("cannot create new session: %v", err)
	}

	if time.Now().Compare(session.ValidUntil) != -1 {
		return nil, errors.New("session is not valid")
	}

	return &session, nil
}

func (s *session) DeleteByUserID(id uint) error {
	return s.db.Where(&models.Session{UserID: id}).Delete(&models.Session{}).Error
}
