package repo

import (
	"errors"

	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Transaction interface {
	Get(userID int) ([]models.Transaction, error)
	Add(t *models.Transaction) (*models.Transaction, error)
}

func NewTransaction(db *gorm.DB, irepo Inventory) Transaction {
	return &transaction{
		db:    db,
		irepo: irepo,
	}
}

type transaction struct {
	irepo Inventory
	db    *gorm.DB
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
	dealMoney := float64(new.Amount) * new.Subject.Price

	err := t.db.Transaction(func(tx *gorm.DB) error {
		sellingItem, err := t.irepo.GetUserInventoryBySecurityID(int(new.Seller.ID), int(new.Subject.SecurityID))
		if err != nil {
			return err
		}

		sellingItem.Amount -= new.Amount
		if _, err := t.irepo.Patch(sellingItem); err != nil {
			return err
		}

		boughtItem, err := t.irepo.GetUserInventoryBySecurityID(int(new.Buyer.ID), int(new.Subject.SecurityID))

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if boughtItem == nil {
			boughtItem.SecurityID = new.Subject.SecurityID
			boughtItem.OwnerID = new.BuyerID
		}

		boughtItem.Amount += new.Amount
		if err := t.irepo.Save(boughtItem); err != nil {
			return err
		}

		if err := t.db.Create(new).Error; err != nil {
			return err
		}

		new.Buyer.Balance -= int(dealMoney)
		if err := tx.Save(new.Buyer).Error; err != nil {
			return err
		}

		new.Seller.Balance += int(dealMoney)
		if err := tx.Save(new.Seller).Error; err != nil {
			return err
		}

		new.Subject.Amount -= new.Amount
		if new.Subject.Amount == 0 {
			new.Subject.Active = false
		}

		if err := tx.Save(new.Subject).Error; err != nil {
			return err
		}

		return nil
	})

	return new, err
}
