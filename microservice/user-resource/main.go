package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/108356037/v1/user-resource-svc/global"
	"github.com/108356037/v1/user-resource-svc/internal/kubeclient"
	"github.com/108356037/v1/user-resource-svc/internal/router"
	"github.com/108356037/v1/user-resource-svc/mq"
	"github.com/108356037/v1/user-resource-svc/pkg"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initSetting() {
	setting, err := pkg.NewSettings()
	Fatal(err)

	err = setting.ReadSection("Server", &global.ServerSetting)
	Fatal(err)

	err = setting.ReadSection("Helm", &global.HelmSetting)
	Fatal(err)
}

func initComponent() {
	err := kubeclient.Init()
	Fatal(err)

	if os.Getenv("RUN_MODE") == "" {
		gin.SetMode(global.ServerSetting.RunMode)
		os.Setenv("MAX_WAIT", "2")
		os.Setenv("BOOTSTRAP_SERVER", "192.168.99.202:9094")
		os.Setenv("SUBSCRIBE_TOPIC", "event")
		os.Setenv("PUBLISH_TOPIC", "event")
		os.Setenv("CONSUMER_GROUPID", "user-resource-svc-group")
	} else {
		gin.SetMode(os.Getenv("RUN_MODE"))
	}

	mq.InitPub()
	mq.InitSub()
}

func main() {
	initSetting()
	initComponent()

	r := router.NewRouter()

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + global.ServerSetting.HttpPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(global.ServerSetting.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(global.ServerSetting.ReadTimeout) * time.Second,
	}
	// go func() {
	// 	mq.RunSubscribe(mq.SimpleLogHandler)
	// }()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Infof("Server started successfully on port %s!", srv.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		if err := mq.KafkaConsumer.Close(); err != nil {
			log.Fatal(err)
		}
		log.Info("Closed Kafka consumer")

		if err := mq.KafkaProducer.Close(); err != nil {
			log.Fatal(err)
		}
		log.Info("Closed Kafka producer")
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Info(err)
	}

	log.Println("Server gracefully shutting down ...")
}
