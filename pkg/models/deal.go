package models

type Deal struct {
	ID         uint `gorm:"primaryKey,autoIncrement"`
	OwnerID    uint
	Owner      User
	SecurityID uint
	Security   Security
	Amount     uint
	Price      float64
	Active     bool
	Sell       bool
}
