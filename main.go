package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"tradeengine/service/db"
)

func main() {
	// init DB
	rootContext := context.Background()
	db.NewDBManager(rootContext)
	// init service

	// Set up signal handling to capture SIGINT and SIGTERM signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigCh:
		os.Exit(1)
	}
}
