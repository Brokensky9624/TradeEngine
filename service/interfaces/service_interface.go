package interfaces

import (
	param "tradeengine/server/web/rest/param"
	memberTypes "tradeengine/service/member/types"
	orderTypes "tradeengine/service/order/types"
	stockTypes "tradeengine/service/stock/types"
	stockInfoTypes "tradeengine/service/stockinfo/types"
	walletTypes "tradeengine/service/wallet/types"
)

type IServiceManager interface {
	MemberService() IMemberSrv
	OrderService() IOrderSrv
	StockInfoService() IStockInfoSrv
	StockService() IStockSrv
}

type IMemberSrv interface {
	Auth(param *param.MemberAuthParam) error
	AuthAndMember(param *param.MemberAuthParam) (*memberTypes.Member, error)
	Create(param param.MemberCreateParam) error
	Edit(param param.MemberEditParam) error
	Delete(param param.MemberDeleteParam) error
	Member(param param.MemberInfoParam) (*memberTypes.Member, error)
	Members() ([]memberTypes.Member, error)
}

type IOrderSrv interface {
	Create(param param.OrderCreateParam) error
	Edit(param param.OrderEditParam) error
	Delete(param param.OrderDeleteParam) error
	OrderInfo(param param.OrderInfoParam) (*orderTypes.Order, error)
	OrderInfoList(param param.OrderInfoListParam) ([]orderTypes.Order, error)
}

type IStockInfoSrv interface {
	Create(param param.StockInfoCreateParam) error
	StockInfo(param param.StockInfoParam) (*stockInfoTypes.StockInfo, error)
	StockInfoList(param param.StockInfoListParam) ([]stockInfoTypes.StockInfo, error)
}

type IWalletSrv interface {
	Create(param param.WalletCreateParam) error
	Edit(param param.WalletEditParam) error
	WalletInfo(param param.WalletInfoParam) (*walletTypes.Wallet, error)
	WalletInfoList(param param.WalletInfoListParam) ([]walletTypes.Wallet, error)
}

type IStockSrv interface {
	Create(param param.OneStockCreateParam) error
	Edit(param param.OneStockEditParam) error
	Delete(param param.OneStockDeleteParam) error
	OneStockInfo(param param.OneStockInfoParam) (*stockTypes.Stock, error)
	OneStockInfoList(param param.OneStockInfoListParam) ([]stockTypes.Stock, error)
}
