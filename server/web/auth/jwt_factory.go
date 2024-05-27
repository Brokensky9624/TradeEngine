package auth

import "github.com/gin-gonic/gin"

type IJWTAuth interface {
	GetRegisterHandler() gin.HandlerFunc
	GetLoginHandler() gin.HandlerFunc
	GetRefreshHandler() gin.HandlerFunc
	GetMiddleWareHandler() gin.HandlerFunc
}

type IJWTAuthFactory interface {
	GetJWTAuth() IJWTAuth
}

type AppleBoyJWTAuthFactory struct {
	identityKey      string
	identityUsername string
	secretKey        string
}

func (f *AppleBoyJWTAuthFactory) initialize() IJWTAuthFactory {
	return nil
}

func (f *AppleBoyJWTAuthFactory) GetJWTAuth() IJWTAuth {
	return nil
}

func GetJWTAuthFactory() IJWTAuthFactory {
	return new(AppleBoyJWTAuthFactory).initialize()
}
