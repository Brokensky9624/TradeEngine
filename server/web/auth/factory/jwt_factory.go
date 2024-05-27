package factory

import (
	authInterfaces "tradeengine/server/web/auth/interfaces"
	serviceInterfaces "tradeengine/service/interfaces"
)

type AppleBoyJWTAuthFactory struct {
	jwtAuth authInterfaces.IJWTAuth
	srvMngr serviceInterfaces.IServiceManager
}

func (f *AppleBoyJWTAuthFactory) InitJWTAuth() {
	f.jwtAuth = NewAppleBoyJWTAuth(f.srvMngr)
}

func (f *AppleBoyJWTAuthFactory) GetJWTAuth() authInterfaces.IJWTAuth {
	return f.jwtAuth
}

func (f *AppleBoyJWTAuthFactory) SetSrvMngr(srvMngr serviceInterfaces.IServiceManager) {
	f.srvMngr = srvMngr
}

func GetJWTAuthFactory() authInterfaces.IJWTAuthFactory {
	return new(AppleBoyJWTAuthFactory)
}
