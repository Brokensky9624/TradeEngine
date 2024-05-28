package factory

import (
	authInterfaces "tradeengine/server/web/auth/interfaces"
	serviceInterfaces "tradeengine/service/interfaces"
)

type AppleBoyJWTAuthFactory struct {
	jwtAuth authInterfaces.IJWTAuth
	srvMngr serviceInterfaces.IServiceManager
}

func NewAppleBoyJWTAuthFactory(srvMngr serviceInterfaces.IServiceManager) *AppleBoyJWTAuthFactory {
	return &AppleBoyJWTAuthFactory{
		srvMngr: srvMngr,
	}
}

func (f *AppleBoyJWTAuthFactory) GetJWTAuth() authInterfaces.IJWTAuth {
	if f.jwtAuth == nil {
		f.jwtAuth = NewAppleBoyJWTAuth(f.srvMngr)
	}
	return f.jwtAuth
}

func GetJWTAuthFactory(srvMngr serviceInterfaces.IServiceManager) authInterfaces.IJWTAuthFactory {
	return NewAppleBoyJWTAuthFactory(srvMngr)
}
