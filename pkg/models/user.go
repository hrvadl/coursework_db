package models

const (
	StockRole   = "stock"
	EmitentRole = "emitent"
)

type User struct {
	ID        uint            `gorm:"primaryKey,autoIncrement"`
	LastName  string          `gorm:"not null"`
	Role      string          `gorm:"not null"`
	FirstName string          `gorm:"not null"`
	Email     string          `gorm:"not null,unique"`
	Password  string          `gorm:"not null"`
	Balance   int             `gorm:"not null"`
	Deals     []Deal          `gorm:"foreignKey:OwnerID"`
	Inventory []InventoryItem `gorm:"foreignKey:OwnerID"`
}
