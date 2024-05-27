package manager

import (
	"sync"

	memberTypes "tradeengine/service/member/types"
)

var (
	manager *ServiceManager
	once    sync.Once
)

func NewManager() *ServiceManager {
	once.Do(func() {
		manager = &ServiceManager{}
	})
	return manager
}

func Manager() *ServiceManager {
	return manager
}

type ServiceManager struct {
	MemberSrv memberTypes.IMemberSrv
}

func (m *ServiceManager) SetMemberService(srv memberTypes.IMemberSrv) *ServiceManager {
	m.MemberSrv = srv
	return m
}

func (m *ServiceManager) MemberService() memberTypes.IMemberSrv {
	return m.MemberSrv
}
