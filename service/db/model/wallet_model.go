package model

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	AvailableMoney float64 `gorm:"not null;default:0"`
	PendingMoney   float64 `gorm:"not null;default:0"`
	OwnerID        uint    `gorm:"not null"`
}

func (Wallet) TableName() string {
	return "wallet"
}
