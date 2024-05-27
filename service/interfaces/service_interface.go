package interfaces

import memberTypes "tradeengine/service/member/types"

type IServiceManager interface {
	MemberService() memberTypes.IMemberSrv
}
