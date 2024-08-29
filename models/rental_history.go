package models

import (
	"time"
)

type RentalHistory struct {
	RentalID          uint       `gorm:"primaryKey;autoIncrement" json:"rental_id"`
	UserID            uint       `gorm:"not null" json:"user_id"`
	CarID             uint       `gorm:"not null" json:"car_id"`
	DriverID          *uint      `json:"driver_id"`
	RentalDate        time.Time  `gorm:"not null" json:"rental_date"`
	ReturnDate        *time.Time `json:"return_date"`
	TotalCost         float64    `gorm:"type:numeric(10,2);not null" json:"total_cost"`
	Status            string     `gorm:"not null;check:status IN ('Book', 'Paid', 'Rent', 'Cancel')" json:"status"`
	PackageID         *uint      `json:"package_id"`
	AirportTransfer   bool       `gorm:"default:false" json:"airport_transfer"`
	PickupLocation    string     `json:"pickup_location"`
	DropoffLocation   string     `json:"dropoff_location"`
	ConciergeServices bool       `gorm:"default:false" json:"concierge_services"`
}
