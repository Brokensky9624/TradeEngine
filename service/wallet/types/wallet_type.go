package types

type Wallet struct {
	ID             uint    `json:"id"`
	OwnerID        uint    `json:"ownerID"`
	AvailableMoney float64 `json:"availableMoney"`
	PendingMoney   float64 `json:"pendingMoney"`
}
