package factory

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"tradeengine/server/web/rest/param"
	"tradeengine/server/web/rest/resp"
	serviceInterfaces "tradeengine/service/interfaces"
	memberTypes "tradeengine/service/member/types"
	"tradeengine/utils/logger"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const (
	identityKey      = "id"                               // indetiyKey for JWT claim
	identityUsername = "name"                             // identityUsername for JWT claim
	secretKey        = "5BYrir4vrBMB2oFJVywHFSrvlim6kCRn" // secret key for JWT encrypt
	realm            = "realm"
	sigAlgo          = "HS256"
	timeout          = time.Hour * 24
	maxRefresh       = time.Hour * 24 * 7
)

type AppleBoyJWTAuth struct {
	srvMngr    serviceInterfaces.IServiceManager
	middleWare *jwt.GinJWTMiddleware
}

func NewAppleBoyJWTAuth(srvMngr serviceInterfaces.IServiceManager) *AppleBoyJWTAuth {
	auth := &AppleBoyJWTAuth{
		srvMngr: srvMngr,
	}
	auth.initialize()
	return auth
}

func (a *AppleBoyJWTAuth) initialize() {
	mid, err := a.GetMiddleWare()
	if err != nil {
		logger.SERVER.Error("failed to initialize AppleBoyJWTAuth for web server, error: %s", err)
	}
	a.middleWare = mid
}

func (a *AppleBoyJWTAuth) GetMiddleWare() (*jwt.GinJWTMiddleware, error) {
	// JWT middleware initialization
	mid, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            realm,             // it used to indicate proteced area or resource in server.
		SigningAlgorithm: sigAlgo,           // specific cryptographic algorithm used to sign the JWT tokens, default is `HS256`
		Key:              []byte(secretKey), // secret key to generate or verify JWT tokens
		Timeout:          timeout,           //  expiration time of JWT tokens, after this time, need relogin or refresh token
		MaxRefresh:       maxRefresh,        // expiration time of refresh JWT token by old JWT token, after this time, need relogin
		IdentityKey:      identityKey,       // this key used to store User identify information in JWT token claims
		PayloadFunc: func(data interface{}) jwt.MapClaims { // this function used to generate JWT token claims by User data
			if v, ok := data.(*memberTypes.Member); ok {
				return jwt.MapClaims{
					identityKey:      v.ID,
					identityUsername: v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} { // this function used to extract JWT token claims from restful request and generate User data
			claims := jwt.ExtractClaims(c)
			id, _ := claims[identityKey].(float64)
			username := claims[identityUsername].(string)
			return &memberTypes.Member{
				ID:   uint(id),
				Name: username,
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) { // this function used to authenticate login user
			var loginUser param.MemberAuthParam
			if err := c.ShouldBindJSON(&loginUser); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			matchUser, err := a.srvMngr.MemberService().AuthAndMember(&loginUser)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return matchUser, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool { // TODO: reserved this function used to check user authorization
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			logger.REST.Debug("failed to auth, err: %s", message)
			c.JSON(code, resp.FailRespObj(errors.New(message)))
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	return mid, err
}

func (a *AppleBoyJWTAuth) GetRegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user param.MemberCreateParam
		if err := c.ShouldBindJSON(&user); err != nil {
			logger.REST.Debug("failed to login, err: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// register member
		if err := a.srvMngr.MemberService().Create(user); err != nil {
			logger.REST.Debug("failed to login, err: %s", err)
			c.JSON(http.StatusOK, resp.FailRespObj(err))
			return
		}
		message := fmt.Sprintf("register member %s successful !", user.Account)
		logger.REST.Info(message)
		c.JSON(http.StatusOK, resp.SuccessRespObj(message, nil))
	}
}

func (a *AppleBoyJWTAuth) GetLoginHandler() gin.HandlerFunc {
	if a.middleWare != nil {
		return a.middleWare.LoginHandler
	}
	return nil
}

func (a *AppleBoyJWTAuth) GetRefreshHandler() gin.HandlerFunc {
	if a.middleWare != nil {
		return a.middleWare.RefreshHandler
	}
	return nil
}

func (a *AppleBoyJWTAuth) GetMiddleWareHandler() gin.HandlerFunc {
	if a.middleWare != nil {
		return a.middleWare.MiddlewareFunc()
	}
	return nil
}
