package undo_buy

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	Name       string
	WalletID   uint
	Wallet     Wallet
	SellOrders []SellOrder
	BuyOrders  []BuyOrder
	Stocks     []Stock
}

type Wallet struct {
	gorm.Model
	AvailableMoney float64
	PendingMoney   float64
	OwnerID        uint
}

type SellOrder struct {
	gorm.Model
	Price       float64
	OwnerID     uint
	StockInfoID uint
	StockInfo   StockInfo
	Quantity    uint
}

type BuyOrder struct {
	gorm.Model
	Price    float64
	Quantity uint
	OwnerID  uint
	Member   Member
}

type StockInfo struct {
	gorm.Model
	Name string
}

type Stock struct {
	gorm.Model
	OwnerID           uint
	StockInfoID       uint
	StockInfo         StockInfo
	AvailableQuantity uint
	PendingQuantity   uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
		return
	}

	db.AutoMigrate(&Member{}, &Wallet{}, &SellOrder{}, &BuyOrder{}, &StockInfo{}, &Stock{})

	member := Member{Name: "John Doe"}
	db.Create(&member)

	wallet := Wallet{AvailableMoney: 1000, PendingMoney: 0, OwnerID: member.ID}
	db.Create(&wallet)
	member.WalletID = wallet.ID
	db.Save(&member)

	createOrUpdateBuyOrder := func(orderID uint, price float64, quantity uint) error {
		return db.Transaction(func(tx *gorm.DB) error {
			var buyOrder BuyOrder
			var isNewOrder bool

			if orderID != 0 {
				if err := tx.First(&buyOrder, orderID).Error; err != nil {
					return err
				}
				isNewOrder = false
			} else {
				buyOrder = BuyOrder{Price: price, Quantity: quantity, OwnerID: member.ID}
				isNewOrder = true
			}

			amountDifference := price*float64(quantity) - buyOrder.Price*float64(buyOrder.Quantity)
			if amountDifference > 0 {
				if wallet.AvailableMoney < amountDifference {
					return fmt.Errorf("insufficient available money")
				}
			}

			if isNewOrder {
				if err := tx.Create(&buyOrder).Error; err != nil {
					return err
				}
			} else {
				buyOrder.Price = price
				buyOrder.Quantity = quantity
				if err := tx.Save(&buyOrder).Error; err != nil {
					return err
				}
			}

			wallet.AvailableMoney -= amountDifference
			wallet.PendingMoney += amountDifference

			if err := tx.Save(&wallet).Error; err != nil {
				return err
			}

			return nil
		})
	}

	updateOrDeleteBuyOrder := func(orderID uint, price float64, quantity uint, delete bool) error {
		return db.Transaction(func(tx *gorm.DB) error {
			var buyOrder BuyOrder

			if err := tx.First(&buyOrder, orderID).Error; err != nil {
				return err
			}

			amountDifference := buyOrder.Price*float64(buyOrder.Quantity) - price*float64(quantity)

			if delete {
				amountDifference = buyOrder.Price * float64(buyOrder.Quantity)
				if err := tx.Delete(&buyOrder).Error; err != nil {
					return err
				}
			} else {
				buyOrder.Price = price
				buyOrder.Quantity = quantity
				if err := tx.Save(&buyOrder).Error; err != nil {
					return err
				}
			}

			wallet.AvailableMoney += amountDifference
			wallet.PendingMoney -= amountDifference

			if wallet.PendingMoney < 0 {
				return fmt.Errorf("pending money cannot be negative")
			}

			if err := tx.Save(&wallet).Error; err != nil {
				return err
			}

			return nil
		})
	}

	if err := createOrUpdateBuyOrder(0, 10.0, 50); err != nil {
		fmt.Printf("Error creating buy order: %v\n", err)
		return
	}

	var buyOrder BuyOrder
	if err := db.Where("owner_id = ?", member.ID).First(&buyOrder).Error; err == nil {
		if err := createOrUpdateBuyOrder(buyOrder.ID, 12.0, 40); err != nil {
			fmt.Printf("Error updating buy order: %v\n", err)
			return
		}
	}

	if err := updateOrDeleteBuyOrder(buyOrder.ID, 0, 0, true); err != nil {
		fmt.Printf("Error deleting buy order: %v\n", err)
		return
	}

	var queriedMember Member
	if err := db.Where("name = ?", "John Doe").Preload("Wallet").Preload("Stocks").Preload("BuyOrders").Preload("SellOrders").First(&queriedMember).Error; err != nil {
		fmt.Printf("Error loading member with associations: %v\n", err)
		return
	}

	fmt.Printf("Member: %v, Wallet: %v, Stocks: %v, BuyOrders: %v, SellOrders: %v\n",
		queriedMember.Name, queriedMember.Wallet, queriedMember.Stocks, queriedMember.BuyOrders, queriedMember.SellOrders)
}
