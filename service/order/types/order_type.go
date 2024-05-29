package types

type OrderType uint

const (
	OrderTypeBuy OrderType = iota + 1
	OrderTypeSell
)

type Order struct {
	ID          uint    `json:"id"`
	OrderType   uint    `json:"orderType"`
	Price       float64 `json:"account"`
	OwnerID     uint    `json:"ownerID"`
	Quantity    uint    `json:"password"`
	StockInfoID uint    `json:"StockInfoID"`
}

func IsBuyOrder(a uint) bool {
	return OrderType(a) == OrderTypeBuy
}

func GetOrderTypeStr(a uint) string {
	d := "buy"
	if !IsBuyOrder(a) {
		d = "sell"
	}
	return d
}
