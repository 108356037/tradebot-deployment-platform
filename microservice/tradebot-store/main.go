package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/108356037/v1/tradebot-store/internal/router"
	"github.com/108356037/v1/tradebot-store/kubeclient"
	log "github.com/sirupsen/logrus"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := kubeclient.Init()
	Fatal(err)

	r := router.NewRouter()
	srv := &http.Server{
		Handler: r,
		Addr:    ":" + os.Getenv("SRV_PORT"),
		// Good practice: enforce timeouts for servers you create!
		//WriteTimeout: time.Duration(global.ServerSetting.WriteTimeout) * time.Second,
		//ReadTimeout:  time.Duration(global.ServerSetting.ReadTimeout) * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Infof("Server started successfully on port %s!", srv.Addr)

	// gracefully shutdown block
	quit := make(chan os.Signal, 1)
	<-quit //blocking here
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		if err := srv.Shutdown(ctx); err != nil {
			log.Warn(err)
		}
		cancel()
	}()

	log.Info("Server gracefully shutting down ...")
}
