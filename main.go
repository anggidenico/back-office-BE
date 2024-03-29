package main

import (
	"context"
	"mf-bo-api/config"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	e := router()

	// Configure server
	s := &http.Server{
		Addr: "0.0.0.0:8030",
	}

	// Start server
	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	//cron running
	if config.Envi == "PROD" {
		// go scheduler()
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx := context.Background()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	} else {
		e.Logger.Info("Gracefully shutdown")
	}
}
