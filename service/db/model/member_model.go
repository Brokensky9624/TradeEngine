package model

import (
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	Name       string      `gorm:"size:100;not null"`
	Account    string      `gorm:"size:100;not null"`
	Password   string      `gorm:"size:100;not null"`
	Email      string      `gorm:"size:200;not null"`
	Phone      string      `gorm:"size:50;not null"`
	Wallet     Wallet      `gorm:"foreignKey:OwnerID;references:ID"`
	SellOrders []SellOrder `gorm:"foreignKey:OwnerID;references:ID"`
	BuyOrders  []BuyOrder  `gorm:"foreignKey:OwnerID;references:ID"`
	Stocks     []Stock     `gorm:"foreignKey:OwnerID;references:ID"`
}

func (Member) TableName() string {
	return "member"
}
