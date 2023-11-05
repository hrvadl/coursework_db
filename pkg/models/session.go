package models

import "time"

type Session struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"not null;unique"`
	UserRole   string `gorm:"not null"`
	ValidUntil time.Time
}
