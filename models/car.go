package models

type Car struct {
	CarID             uint    `gorm:"primaryKey;autoIncrement" json:"car_id"`
	Name              string  `gorm:"not null" json:"name"`
	StockAvailability int     `gorm:"not null;check:stock_availability >= 0" json:"stock_availability"`
	RentalCosts       float64 `gorm:"type:numeric(10,2);not null" json:"rental_costs"`
	Category          string  `gorm:"not null" json:"category"`
	Make              string  `gorm:"not null" json:"make"`
	Model             string  `gorm:"not null" json:"model"`
	Transmission      string  `gorm:"not null" json:"transmission"`
	Year              int     `gorm:"not null" json:"year"`
	FuelType          string  `gorm:"not null" json:"fuel_type"`
	Class             string  `gorm:"not null" json:"class"`
}
