package param

import "tradeengine/utils/tool"

type StockInfoCreateParam struct {
	Name string `json:"name" required:"true"`
}

func (param StockInfoCreateParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type StockInfoParam struct {
	ID uint `json:"stockInfoID" required:"true"`
}

func (param StockInfoParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type StockInfoListParam struct {
	OrderBy   string `json:"orderBy"`
	OrderDesc bool   `json:"orderDirect"`
}

func (param StockInfoListParam) Check() error {
	return tool.CheckRequiredFields(param)
}
