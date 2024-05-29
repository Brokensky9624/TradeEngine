package model

import "gorm.io/gorm"

type StockInfo struct {
	gorm.Model
	Name string `gorm:"size:50;not null"`
}

func (StockInfo) TableName() string {
	return "stock_info"
}
