package param

import "tradeengine/utils/tool"

type OneStockCreateParam struct {
	OwnerID           uint `json:"ownerID" required:"true"`
	StockInfoID       uint `json:"stockInfoID" required:"true"`
	AvailableQuantity uint `json:"availableQuantity" required:"true"`
	PendingQuantity   uint `json:"pendingQuantity"`
}

func (param OneStockCreateParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type OneStockEditParam struct {
	ID                uint `json:"stockID" required:"true"`
	OwnerID           uint `json:"ownerID" required:"true"`
	AvailableQuantity uint `json:"availableQuantity"`
	PendingQuantity   uint `json:"pendingQuantity"`
}

func (param OneStockEditParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type OneStockDeleteParam struct {
	ID      uint `json:"stockID" required:"true"`
	OwnerID uint `json:"ownerID" required:"true"`
}

func (param OneStockDeleteParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type OneStockInfoParam struct {
	ID      uint `json:"walletID" required:"true"`
	OwnerID uint `json:"ownerID" required:"true"`
}

func (param OneStockInfoParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type OneStockInfoListParam struct {
	OrderBy   string `json:"orderBy"`
	OrderDesc bool   `json:"orderDirect"`
}

func (param OneStockInfoListParam) Check() error {
	return tool.CheckRequiredFields(param)
}
