package param

import (
	"tradeengine/utils/tool"
)

type OrderCreateParam struct {
	OrderType   uint    `json:"orderType" required:"true"`
	Price       float64 `json:"price" required:"true"`
	OwnerID     uint    `json:"ownerID" required:"true"`
	StockInfoID uint    `json:"StockInfoID" required:"true"`
	Quantity    uint    `json:"quantity" required:"true"`
}

func (param OrderCreateParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type OrderEditParam struct {
	ID        uint    `json:"orderID" required:"true"`
	OrderType uint    `json:"orderType" required:"true"`
	Price     float64 `json:"price"`
	OwnerID   uint    `json:"ownerID"`
	Quantity  uint    `json:"quantity"`
}

func (param OrderEditParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type OrderDeleteParam struct {
	ID        uint `json:"orderID" required:"true"`
	OrderType uint `json:"orderType" required:"true"`
	OwnerID   uint `json:"ownerID"`
}

func (param OrderDeleteParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type OrderInfoParam struct {
	ID        uint `json:"orderID" required:"true"`
	OrderType uint `json:"orderType" required:"true"`
}

func (param OrderInfoParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type OrderInfoListParam struct {
	OrderType uint   `json:"orderType" required:"true"`
	OwnerID   uint   `json:"ownerID"`
	OrderBy   string `json:"orderBy"`
	OrderDesc bool   `json:"orderDirect"`
}

func (param OrderInfoListParam) Check() error {
	return tool.CheckRequiredFields(param)
}
