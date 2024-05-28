package interfaces

import (
	"github.com/gin-gonic/gin"
)

type IJWTAuth interface {
	GetRegisterHandler() gin.HandlerFunc
	GetLoginHandler() gin.HandlerFunc
	GetRefreshHandler() gin.HandlerFunc
	GetMiddleWareHandler() gin.HandlerFunc
}

type IJWTAuthFactory interface {
	GetJWTAuth() IJWTAuth
}
