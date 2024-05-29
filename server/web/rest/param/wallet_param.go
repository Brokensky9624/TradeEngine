package param

import "tradeengine/utils/tool"

type WalletCreateParam struct {
	OwnerID        uint    `json:"ownerID" required:"true"`
	AvailableMoney float64 `json:"availableMoney" required:"true"`
	PendingMoney   float64 `json:"pendingMoney" required:"true"`
}

func (param WalletCreateParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type WalletEditParam struct {
	ID             uint
	OwnerID        uint    `json:"ownerID" required:"true"`
	AvailableMoney float64 `json:"availableMoney"`
	PendingMoney   float64 `json:"pendingMoney"`
}

func (param WalletEditParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type WalletInfoParam struct {
	ID uint `json:"walletID" required:"true"`
}

func (param WalletInfoParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type WalletInfoListParam struct {
	OrderBy   string `json:"orderBy"`
	OrderDesc bool   `json:"orderDirect"`
}

func (param WalletInfoListParam) Check() error {
	return tool.CheckRequiredFields(param)
}
