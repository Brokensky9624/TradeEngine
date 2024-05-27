package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tradeengine/server/web"
	"tradeengine/server/web/rest/member"
	"tradeengine/service/db"
	"tradeengine/utils/logger"
)

func main() {
	logger.SERVER.Info("process has started")

	rootContext, rootCtxCancel := context.WithCancel(context.Background())
	defer func() {
		logger.SERVER.Info("process is closing")
		rootCtxCancel()
		time.Sleep(10 * time.Second)
		logger.SERVER.Info("process has closed")
	}()

	// init DB Manager handling database
	dbMngr := db.NewDBManager(rootContext)
	dbMngr.Run()

	// init web server handling restful
	webSrv := web.NewWebServer(rootContext)
	webSrv.Prepare()
	webSrv.Run()

	// init service
	member.NewService()

	// Set up signal handling to capture SIGINT and SIGTERM signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	logger.SERVER.Info("received signal, shutting down")
	// select {
	// case <-sigCh:
	// 	os.Exit(1)
	// }
}
