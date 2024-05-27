package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
	"tradeengine/server/web/auth"
	"tradeengine/server/web/rest/member"
	"tradeengine/utils/logger"
	"tradeengine/utils/panichandle"

	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
)

const (
	host               = "127.0.0.1"
	port               = "8816"
	readTimeout        = 30 * time.Second
	writeTimeout       = 30 * time.Second
	httpRestartPeriod  = 5 * time.Second
	httpsRestartPeriod = 5 * time.Second
	cMuxRestartPeriod  = 5 * time.Second
)

type WebServer struct {
	*gin.Engine
	cMux        cmux.CMux
	httpServer  http.Server
	httpsServer http.Server
	jwtAuth     auth.IJWTAuth
	mainGroup   *gin.RouterGroup
}

var (
	webServer *WebServer
	once      sync.Once
)

func NewWebServer(ctx context.Context) *WebServer {
	once.Do(func() {
		logger.SERVER.Info("web server initializing")
		webServer = createWebServer()
	})
	return webServer
}

func createWebServer() *WebServer {
	return &WebServer{
		Engine: gin.Default(),
	}
}

func (w *WebServer) Prepare() {
	w.loadJWT()
	w.loadMainGroup()
	w.registerBaseRoute()
	w.loadMiddleWare()
	w.registerRoute()
}

func (w *WebServer) loadJWT() {
	w.jwtAuth = auth.GetJWTAuthFactory().GetJWTAuth()
}

func (w *WebServer) loadMainGroup() {
	w.mainGroup = w.Group("/api")
}

func (w *WebServer) registerBaseRoute() {
	w.POST("register", w.jwtAuth.GetRegisterHandler())
	w.POST("/login", w.jwtAuth.GetLoginHandler())
	w.POST("/refresh-token", w.jwtAuth.GetRefreshHandler())
}

func (w *WebServer) loadMiddleWare() {
	w.mainGroup.Use(w.jwtAuth.GetMiddleWareHandler())
}

func (w *WebServer) registerRoute() {
	member.NewREST(w.mainGroup).RegisterRoute()
}

// support auto restart version
func (w *WebServer) Run() {
	go func() {
		defer panichandle.PanicHandle()
		l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			panic(err)
		}
		defer l.Close()
		w.cMux = cmux.New(l)
		httpL := w.cMux.Match(cmux.HTTP1Fast())
		httpsL := w.cMux.Match(cmux.Any())

		// start cmux
		go func() {
			for {
				if err := w.cMux.Serve(); err != nil {
					fmt.Printf("cmux server error: %s\n", err)
					select {
					case <-w.ctx.Done():
						return
					default:
						time.Sleep(cMuxRestartPeriod)
						w.cMux = cmux.New(l)
						httpL = w.cMux.Match(cmux.HTTP1Fast())
						httpsL = w.cMux.Match(cmux.Any())
					}
				}
			}
		}()

		// start http server
		go func() {
			httpSrv := getHTTPServer(w.Engine)
			for {
				if err := httpSrv.Serve(httpL); err != nil && err != http.ErrServerClosed {
					fmt.Printf("HTTP server error: %s\n", err)
					select {
					case <-w.ctx.Done():
						return
					default:
						time.Sleep(httpRestartPeriod)
					}
				}
			}
		}()

		// start https server
		go func() {
			httpsSrv := getHTTPsServer(w.Engine)
			for {
				if err := httpsSrv.Serve(httpsL); err != nil && err != http.ErrServerClosed {
					fmt.Printf("HTTPS server error: %s\n", err)
					select {
					case <-w.ctx.Done():
						return
					default:
						time.Sleep(httpsRestartPeriod)
					}
				}
			}
		}()

		// close all web server
		select {
		case <-w.ctx.Done():
			logger.SERVER.Info("web servers are shutting down")
			w.cMux.Close()
			httpL.Close()
			httpsL.Close()
			logger.SERVER.Info("all web servers were closed")

		}
	}()
}

func Server() *WebServer {
	return webServer
}

func getHTTPServer(router *gin.Engine) *http.Server {
	return &http.Server{
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func getHTTPsServer(router *gin.Engine) *http.Server {
	server := getHTTPServer(router)
	server.TLSConfig = getTLSConfig()
	return server
}
