package models

type Deal struct {
	ID         uint  `gorm:"primaryKey,autoIncrement"`
	OwnerID    uint  `gorm:"index:idx_owner"`
	Owner      *User `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	SecurityID uint
	Security   *Security `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Amount     uint      `gorm:"check:amount > 0"`
	Price      float64   `gorm:"check:price > 0"`
	Active     bool
	Sell       bool
}
