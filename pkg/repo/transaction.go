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

// TODO: delete when 0. AJAX request
func (t *transaction) Add(new *models.Transaction) (*models.Transaction, error) {
	dealMoney := float64(new.Amount) * new.Subject.Price

	// TODO: why this does not work?
	if err := t.db.Create(new).Error; err != nil {
		return nil, err
	}

	err := t.db.Transaction(func(tx *gorm.DB) error {
		new.Buyer.Balance -= int(dealMoney)
		if err := tx.Save(new.Buyer).Error; err != nil {
			return err
		}

		new.Seller.Balance += int(dealMoney)
		if err := tx.Save(new.Seller).Error; err != nil {
			return err
		}

		new.Subject.Amount -= new.Amount
		if err := tx.Save(new.Subject).Error; err != nil {
			return err
		}

		return nil
	})

	return new, err
}
