package models

type Transaction struct {
	ID uint `gorm:"primaryKey,autoIncrement"`

	BuyerID uint
	Buyer   User

	SellerID uint
	Seller   User

	SubjectID uint
	Subject   Deal

	Amount uint
}
