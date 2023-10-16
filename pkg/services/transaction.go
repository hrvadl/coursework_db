package services

import (
	"errors"

	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type Transaction interface {
	Get(userID int) ([]models.Transaction, error)
	Add(t *models.Transaction) (*models.Transaction, error)
}

func NewTransaction(tr repo.Transaction) Transaction {
	return &transaction{
		tr: tr,
	}
}

type transaction struct {
	tr repo.Transaction
}

func (t *transaction) Get(userID int) ([]models.Transaction, error) {
	return t.tr.Get(userID)
}

func (t *transaction) Add(tr *models.Transaction) (*models.Transaction, error) {
	if tr.BuyerID == 0 {
		return nil, errors.New("buyer ID cannot be empty")
	}

	if tr.SellerID == 0 {
		return nil, errors.New("seller ID cannot be empty")
	}

	if tr.SubjectID == 0 {
		return nil, errors.New("subject ID cannot be empty")
	}

	return t.tr.Add(tr)
}
