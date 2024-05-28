package order

import (
	"gorm.io/gorm"
)

type OrderBuy struct {
	gorm.Model
	Price    uint `gorm:"type:int unsigned;check:price <= 1000"`
	OwnerID  uint `gorm:"type:int unsigned;check:owner_id <= 1000"`
	Quantity uint `gorm:"type:int unsigned;check:quantity <= 1000"`
	StockID  uint `gorm:"type:int unsigned;check:stock_id <= 1000"`
	Status   uint `gorm:"type:int unsigned;check:status <= 1000"`
}

func (OrderBuy) TableName() string {
	return "order_buy"
}
