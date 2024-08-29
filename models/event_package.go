package models

type EventPackage struct {
	PackageID   uint    `gorm:"primaryKey;autoIncrement" json:"package_id"`
	PackageName string  `gorm:"not null" json:"package_name"`
	Description string  `json:"description"`
	Cost        float64 `gorm:"type:numeric(10,2);not null" json:"cost"`
}
