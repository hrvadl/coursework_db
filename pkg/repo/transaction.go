package repo

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Transaction interface {
	Get(userID int) ([]models.Transaction, error)
	Add(t *models.Transaction) (*models.Transaction, error)
}

func NewTransaction(db *gorm.DB) Transaction {
	return &transaction{db}
}

type transaction struct {
	db *gorm.DB
}

func (t *transaction) Get(userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	res := t.db.
		Model(&models.Transaction{}).
		Preload(clause.Associations).
		Preload("Subject.Security").
		Where(&models.Transaction{SellerID: uint(userID)}).
		Or(&models.Transaction{BuyerID: uint(userID)}).
		Find(&transactions)

	if err := res.Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *transaction) Add(new *models.Transaction) (*models.Transaction, error) {
	if err := t.db.Create(new).Error; err != nil {
		return nil, err
	}

	return new, nil
}
