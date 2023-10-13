package models

type Deal struct {
	ID         uint `gorm:"primaryKey,autoIncrement"`
	OwnerID    uint
	SecurityID uint
	Amount     uint
	Price      float64
	Active     bool
}
