package types

type Stock struct {
	ID                uint `json:"id"`
	OwnerID           uint `json:"ownerID"`
	StockInfoID       uint `json:"stockInfoID"`
	AvailableQuantity uint `json:"availableQuantity"`
	PendingQuantity   uint `json:"pendingQuantity"`
}
