package manager

import (
	"sync"

	serverInterfaces "tradeengine/service/interfaces"
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
	MemberSrv serverInterfaces.IMemberSrv
	OrderSrv  serverInterfaces.IOrderSrv
}

func (m *ServiceManager) SetMemberService(srv serverInterfaces.IMemberSrv) *ServiceManager {
	m.MemberSrv = srv
	return m
}

func (m *ServiceManager) MemberService() serverInterfaces.IMemberSrv {
	return m.MemberSrv
}

func (m *ServiceManager) SetOrderService(srv serverInterfaces.IOrderSrv) *ServiceManager {
	m.OrderSrv = srv
	return m
}

func (m *ServiceManager) OrderService() serverInterfaces.IOrderSrv {
	return m.OrderSrv
}
