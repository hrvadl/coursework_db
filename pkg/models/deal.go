package models

type Deal struct {
	ID         uint `gorm:"primaryKey,autoIncrement"`
	OwnerID    uint
	Owner      *User `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	SecurityID uint
	Security   *Security `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Amount     uint
	Price      float64
	Active     bool
	Sell       bool
}
