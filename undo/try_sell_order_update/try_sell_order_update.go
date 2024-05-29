package undo_sell

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	Name       string      `gorm:"size:100;not null"`
	Wallet     Wallet      `gorm:"foreignKey:OwnerID;references:ID"`
	SellOrders []SellOrder `gorm:"foreignKey:OwnerID;references:ID"`
	BuyOrders  []BuyOrder  `gorm:"foreignKey:OwnerID;references:ID"`
	Stocks     []Stock     `gorm:"foreignKey:OwnerID;references:ID"`
}

type Wallet struct {
	gorm.Model
	AvailableMoney float64 `gorm:"not null;default:0"`
	PendingMoney   float64 `gorm:"not null;default:0"`
	OwnerID        uint    `gorm:"not null"`
}

type SellOrder struct {
	gorm.Model
	Price       float64   `gorm:"not null;check:price <= 1000"`
	OwnerID     uint      `gorm:"not null"`
	StockInfoID uint      `gorm:"not null"`
	StockInfo   StockInfo `gorm:"foreignKey:StockInfoID;references:ID"`
	Member      Member    `gorm:"foreignKey:OwnerID;references:ID"`
	Quantity    uint      `gorm:"not null;check:quantity <= 1000"`
}

type BuyOrder struct {
	gorm.Model
	Price       float64   `gorm:"not null;check:price <= 1000"`
	Quantity    uint      `gorm:"not null;check:quantity <= 1000"`
	OwnerID     uint      `gorm:"not null"`
	StockInfoID uint      `gorm:"not null"`
	StockInfo   StockInfo `gorm:"foreignKey:StockInfoID;references:ID"`
	Member      Member    `gorm:"foreignKey:OwnerID;references:ID"`
}

type StockInfo struct {
	gorm.Model
	Name string `gorm:"size:50;not null"`
}

type Stock struct {
	gorm.Model
	OwnerID           uint      `gorm:"not null"`
	StockInfoID       uint      `gorm:"not null"`
	StockInfo         StockInfo `gorm:"foreignKey:StockInfoID;references:ID"`
	Member            Member    `gorm:"foreignKey:OwnerID;references:ID"`
	AvailableQuantity uint      `gorm:"not null;check:available_quantity <= 1000;default:0"`
	PendingQuantity   uint      `gorm:"not null;check:pending_quantity <= 1000;default:0"`
}

func main() {
	// 初始化数据库连接
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
		return
	}

	// 自动迁移模式
	db.AutoMigrate(&Member{}, &Wallet{}, &SellOrder{}, &BuyOrder{}, &StockInfo{}, &Stock{})

	// 假设 stock id 为 14
	stockID := uint(14)

	err = db.Transaction(func(tx *gorm.DB) error {
		var sellOrder SellOrder
		var buyOrder BuyOrder

		// 找到所有 stock id 等于 14 且 price 最低且取更新時間最早的一笔 sell order
		if err := tx.Where("stock_info_id = ?", stockID).Order("price asc, updated_at asc").First(&sellOrder).Error; err != nil {
			return fmt.Errorf("no suitable sell order found")
		}

		// 找到所有 stock id 等于 14 且 price 取最高且更新時間最早的一笔 buy order
		if err := tx.Where("stock_info_id = ?", stockID).Order("price desc, updated_at asc").First(&buyOrder).Error; err != nil {
			return fmt.Errorf("no suitable buy order found")
		}

		// 检查交易是否成功
		if buyOrder.Price >= sellOrder.Price {
			// 获取买卖双方会员
			var sellMember, buyMember Member
			if err := tx.First(&sellMember, sellOrder.OwnerID).Error; err != nil {
				return err
			}
			if err := tx.First(&buyMember, buyOrder.OwnerID).Error; err != nil {
				return err
			}

			// 获取买卖双方钱包
			var sellWallet, buyWallet Wallet
			if err := tx.Where("owner_id = ?", sellMember.ID).First(&sellWallet).Error; err != nil {
				return err
			}
			if err := tx.Where("owner_id = ?", buyMember.ID).First(&buyWallet).Error; err != nil {
				return err
			}

			transactionAmount := float64(sellOrder.Quantity) * sellOrder.Price

			// 交易成功
			if buyOrder.Quantity > sellOrder.Quantity {
				// buy order quantity > sell order quantity
				// 更新 buy order 的 quantity
				if err := tx.Model(&buyOrder).Update("quantity", gorm.Expr("quantity - ?", sellOrder.Quantity)).Error; err != nil {
					return err
				}
				// 删除 sell order
				if err := tx.Delete(&sellOrder).Error; err != nil {
					return err
				}
			} else if buyOrder.Quantity == sellOrder.Quantity {
				// buy order quantity == sell order quantity
				// 删除 buy order 和 sell order
				if err := tx.Delete(&buyOrder).Error; err != nil {
					return err
				}
				if err := tx.Delete(&sellOrder).Error; err != nil {
					return err
				}
			} else {
				// buy order quantity < sell order quantity
				// 更新 sell order 的 quantity
				if err := tx.Model(&sellOrder).Update("quantity", gorm.Expr("quantity - ?", buyOrder.Quantity)).Error; err != nil {
					return err
				}
				// 删除 buy order
				if err := tx.Delete(&buyOrder).Error; err != nil {
					return err
				}
			}

			// 更新钱包和库存
			sellWallet.AvailableMoney += transactionAmount
			if err := tx.Model(&Stock{}).Where("id = ?", sellOrder.StockInfoID).Update("pending_quantity", gorm.Expr("pending_quantity - ?", sellOrder.Quantity)).Error; err != nil {
				return err
			}

			buyWallet.PendingMoney -= transactionAmount
			if err := tx.Model(&Stock{}).Where("id = ?", buyOrder.StockInfoID).Update("available_quantity", gorm.Expr("available_quantity + ?", buyOrder.Quantity)).Error; err != nil {
				return err
			}

			if err := tx.Save(&sellWallet).Error; err != nil {
				return err
			}
			if err := tx.Save(&buyWallet).Error; err != nil {
				return err
			}
		} else {
			return fmt.Errorf("buy order price is less than sell order price")
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Transaction failed: %v\n", err)
	} else {
		fmt.Println("Transaction succeeded")
	}
}
