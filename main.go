package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"tradeengine/service/db"
	"tradeengine/utils/logger"
)

func main() {
	logger.SERVER.Info("Process has started")

	// init DB Manager handling database
	rootContext, rootCtxCancel := context.WithCancel(context.Background())
	dbMngr := db.NewDBManager(rootContext)
	dbMngrFinished := dbMngr.Run()
	defer func() {
		rootCtxCancel()
		<-dbMngrFinished // wait until all DB closed
		logger.SERVER.Info("All DB has closed")
		logger.SERVER.Info("Process has closed")
	}()

	// init web server handling restful

	// init service

	// Set up signal handling to capture SIGINT and SIGTERM signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	logger.SERVER.Info("Received signal, shutting down...")
	// select {
	// case <-sigCh:
	// 	os.Exit(1)
	// }
}
