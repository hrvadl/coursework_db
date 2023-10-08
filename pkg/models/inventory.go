package models

type InventoryItem struct {
	ID         uint `gorm:"primaryKey,autoIncrement"`
	SecurityID uint

	Security Security
	OwnerID  uint

	Owner  User
	Amount uint
}
