package model

import "gorm.io/gorm"

type Stock struct {
	gorm.Model
	OwnerID           uint      `gorm:"not null"`
	StockInfoID       uint      `gorm:"not null"`
	StockInfo         StockInfo `gorm:"foreignKey:StockInfoID;references:ID"`
	Member            Member    `gorm:"foreignKey:OwnerID;references:ID"`
	AvailableQuantity uint      `gorm:"not null;check:available_quantity <= 1000;default:0"`
	PendingQuantity   uint      `gorm:"not null;check:pending_quantity <= 1000;default:0"`
}

func (Stock) TableName() string {
	return "stock"
}
