package repo

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
)

type Session interface {
	Create(u *models.User) (*models.Session, error)
	GetByID(id uint) (*models.Session, error)
}

func NewSession(db *gorm.DB) Session {
	return &session{
		db: db,
	}
}

type session struct {
	db *gorm.DB
}

func (s *session) Create(u *models.User) (*models.Session, error) {
	session := &models.Session{
		UserID:     u.ID,
		UserRole:   u.Role,
		ValidUntil: time.Now().Add(time.Hour * 24),
	}

	_ = s.DeleteByUserID(u.ID)
	res := s.db.Create(session)

	if res.Error != nil {
		return nil, res.Error
	}

	return session, nil
}

func (s *session) GetByID(id uint) (*models.Session, error) {
	var session models.Session

	if err := s.db.First(&session, id).Error; err != nil {
		return nil, fmt.Errorf("cannot find session with ID %v: %v", id, err)
	}

	log.Printf("now: %v, valid till: %v", time.Now(), session.ValidUntil)
	if time.Now().Compare(session.ValidUntil) != -1 {
		return nil, errors.New("session is not valid")
	}

	return &session, nil
}

func (s *session) DeleteByUserID(id uint) error {
	return s.db.Where(&models.Session{UserID: id}).Delete(&models.Session{}).Error
}
