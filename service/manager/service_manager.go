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
	MemberSrv    serverInterfaces.IMemberSrv
	OrderSrv     serverInterfaces.IOrderSrv
	StockInfoSrv serverInterfaces.IStockInfoSrv
	WalletSrv    serverInterfaces.IWalletSrv
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

func (m *ServiceManager) SetStockInfoService(srv serverInterfaces.IStockInfoSrv) *ServiceManager {
	m.StockInfoSrv = srv
	return m
}

func (m *ServiceManager) StockInfoService() serverInterfaces.IStockInfoSrv {
	return m.StockInfoSrv
}

func (m *ServiceManager) SetWalletService(srv serverInterfaces.IWalletSrv) *ServiceManager {
	m.WalletSrv = srv
	return m
}

func (m *ServiceManager) WalletService() serverInterfaces.IWalletSrv {
	return m.WalletSrv
}
