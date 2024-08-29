package models

type Membership struct {
	MembershipID  uint   `gorm:"primaryKey;autoIncrement" json:"membership_id"`
	UserID        uint   `gorm:"not null" json:"user_id"`
	DiscountLevel string `gorm:"not null;check:discount_level IN ('Silver', 'Gold', 'Platinum')" json:"discount_level"`
}
