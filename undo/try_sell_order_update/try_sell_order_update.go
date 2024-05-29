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
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
		return
	}

	db.AutoMigrate(&Member{}, &Wallet{}, &SellOrder{}, &BuyOrder{}, &StockInfo{}, &Stock{})

	stockID := uint(14)

	err = db.Transaction(func(tx *gorm.DB) error {
		var sellOrder SellOrder
		var buyOrder BuyOrder

		if err := tx.Where("stock_info_id = ?", stockID).Order("price asc, updated_at asc").First(&sellOrder).Error; err != nil {
			return fmt.Errorf("no suitable sell order found")
		}

		if err := tx.Where("stock_info_id = ?", stockID).Order("price desc, updated_at asc").First(&buyOrder).Error; err != nil {
			return fmt.Errorf("no suitable buy order found")
		}

		if buyOrder.Price >= sellOrder.Price {
			var sellMember, buyMember Member
			if err := tx.First(&sellMember, sellOrder.OwnerID).Error; err != nil {
				return err
			}
			if err := tx.First(&buyMember, buyOrder.OwnerID).Error; err != nil {
				return err
			}

			var sellWallet, buyWallet Wallet
			if err := tx.Where("owner_id = ?", sellMember.ID).First(&sellWallet).Error; err != nil {
				return err
			}
			if err := tx.Where("owner_id = ?", buyMember.ID).First(&buyWallet).Error; err != nil {
				return err
			}

			transactionAmount := float64(sellOrder.Quantity) * sellOrder.Price

			if buyOrder.Quantity > sellOrder.Quantity {
				if err := tx.Model(&buyOrder).Update("quantity", gorm.Expr("quantity - ?", sellOrder.Quantity)).Error; err != nil {
					return err
				}
				if err := tx.Delete(&sellOrder).Error; err != nil {
					return err
				}
			} else if buyOrder.Quantity == sellOrder.Quantity {
				if err := tx.Delete(&buyOrder).Error; err != nil {
					return err
				}
				if err := tx.Delete(&sellOrder).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(&sellOrder).Update("quantity", gorm.Expr("quantity - ?", buyOrder.Quantity)).Error; err != nil {
					return err
				}
				if err := tx.Delete(&buyOrder).Error; err != nil {
					return err
				}
			}

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
