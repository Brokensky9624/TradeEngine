package model

import (
	"gorm.io/gorm"
)

type BuyOrder struct {
	gorm.Model
	Price       float64   `gorm:"not null;check:price <= 1000"`
	Quantity    uint      `gorm:"not null;check:quantity <= 1000"`
	OwnerID     uint      `gorm:"not null"`
	StockInfoID uint      `gorm:"not null"`
	StockInfo   StockInfo `gorm:"foreignKey:StockInfoID;references:ID"`
	Member      Member    `gorm:"foreignKey:OwnerID;references:ID"`
}

func (BuyOrder) TableName() string {
	return "buy_order"
}
