package interfaces

import (
	"tradeengine/service/interfaces"

	"github.com/gin-gonic/gin"
)

type IJWTAuth interface {
	GetRegisterHandler() gin.HandlerFunc
	GetLoginHandler() gin.HandlerFunc
	GetRefreshHandler() gin.HandlerFunc
	GetMiddleWareHandler() gin.HandlerFunc
}

type IJWTAuthFactory interface {
	InitJWTAuth()
	GetJWTAuth() IJWTAuth
	SetSrvMngr(mngr interfaces.IServiceManager)
}
