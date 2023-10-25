package services

import (
	"errors"

	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type User interface {
	GetByID(id int) (*models.User, error)
	Get() ([]models.User, error)
	Patch(user *models.User) (*models.User, error)
	WithdrawMoney(userID, amount int) error
	AddMoney(userID, amount int) error
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

func (s *stock) WithdrawMoney(userID, amount int) error {
	profile, err := s.GetByID(userID)

	if err != nil {
		return err
	}

	if amount < 1 {
		return errors.New("amount must be greater than zero")
	}

	if profile.Balance-amount < 0 {
		return errors.New("total amount must be greater than zero")
	}

	profile.Balance -= amount
	_, err = s.Patch(profile)
	return err
}

func (s *stock) AddMoney(userID, amount int) error {
	profile, err := s.GetByID(userID)

	if err != nil {
		return err
	}

	if amount < 1 {
		return errors.New("amount must be greater than zero")
	}

	profile.Balance += amount
	_, err = s.Patch(profile)
	return err
}
