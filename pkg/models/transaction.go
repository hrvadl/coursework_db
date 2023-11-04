package models

type Transaction struct {
	ID      uint  `gorm:"primaryKey,autoIncrement"`
	Buyer   *User `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	BuyerID uint

	Seller   *User `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	SellerID uint

	Subject   *Deal `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	SubjectID uint
	Amount    uint
}
