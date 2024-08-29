package models

type Driver struct {
	DriverID        uint    `gorm:"primaryKey;autoIncrement" json:"driver_id"`
	Name            string  `gorm:"not null" json:"name"`
	PhoneNumber     string  `gorm:"not null" json:"phone_number"`
	LicenseNumber   string  `gorm:"unique;not null" json:"license_number"`
	ExperienceYears int     `gorm:"not null" json:"experience_years"`
	Rating          float64 `gorm:"type:numeric(3,2);default:5.0" json:"rating"`
}
