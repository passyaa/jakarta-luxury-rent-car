package models

type User struct {
	UserID        uint    `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Email         string  `gorm:"unique;not null" json:"email"`
	Password      string  `gorm:"not null" json:"password"`
	PhoneNumber   string  `gorm:"type:varchar(15);not null" json:"phone_number"`
	Address       string  `gorm:"not null" json:"address"`
	DepositAmount float64 `gorm:"type:numeric(10,2);default:0" json:"deposit_amount"`
	Role          string  `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
}
