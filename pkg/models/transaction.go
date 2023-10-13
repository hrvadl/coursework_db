package models

type Transaction struct {
	ID      uint `gorm:"primaryKey,autoIncrement"`
	Buyer   User
	BuyerID uint

	Seller   User
	SellerID uint

	SubjectID uint
	Amount    uint
}
