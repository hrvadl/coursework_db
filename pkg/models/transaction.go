package models

type Transaction struct {
	ID      uint `gorm:"primaryKey,autoIncrement"`
	Buyer   *User
	BuyerID uint

	Seller   *User
	SellerID uint

	Subject   *Deal
	SubjectID uint
	Amount    uint
}
