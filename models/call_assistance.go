package models

import (
	"time"
)

type CallAssistance struct {
	CallAssistanceID   uint      `gorm:"primaryKey;autoIncrement" json:"assistance_id"`
	RentalID       uint      `gorm:"not null" json:"rental_id"`
	UserID         uint      `gorm:"not null" json:"user_id"` // New field to identify the user
	CallAssistanceDate time.Time `gorm:"not null" json:"callassistance_date"`
	Description    string    `gorm:"not null" json:"description"`
	Location       string    `gorm:"not null" json:"location"` // New field for assistance location
}
