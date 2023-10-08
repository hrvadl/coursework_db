package models

import "time"

type Session struct {
	UserID     uint      `gorm:"primaryKey"`
	ValidUntil time.Time `gorm:"primaryKey"`
}
