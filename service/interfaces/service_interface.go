package interfaces

import (
	param "tradeengine/server/web/rest/param"
	memberTypes "tradeengine/service/member/types"
	orderTypes "tradeengine/service/order/types"
)

type IServiceManager interface {
	MemberService() IMemberSrv
	OrderService() IOrderSrv
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
