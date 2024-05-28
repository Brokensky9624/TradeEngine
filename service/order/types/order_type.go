package types

type OrderType uint

const (
	OrderTypeBuy OrderType = iota
	OrderTypeSell
)

type StatusType uint

const (
	OrderStatusNew StatusType = iota
)

type Order struct {
	ID        uint `json:"id"`
	OrderType uint `json:"orderType"`
	Price     uint `json:"account"`
	OwnerID   uint `json:"ownerID"`
	Quantity  uint `json:"password"`
	StockID   uint `json:"stockID"`
	Status    uint `json:"status"`
}

func IsOrderBuy(a uint) bool {
	return OrderType(a) == OrderTypeBuy
}

func GetOrderTypeStr(a uint) string {
	d := "buy"
	if !IsOrderBuy(a) {
		d = "sell"
	}
	return d
}
