package models

type Security struct {
	ID   uint   `gorm:"primaryKey,autoIncrement"`
	Name string `gorm:"not null"`
}
