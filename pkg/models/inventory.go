package models

type InventoryItem struct {
	ID         uint `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	SecurityID uint
	Security   Security `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	OwnerID    uint
	Amount     uint
}
