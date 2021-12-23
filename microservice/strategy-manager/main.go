package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/108356037/v1/strategy-manager/internal/database/mongo"
	"github.com/108356037/v1/strategy-manager/internal/grpcclient"

	"github.com/gin-gonic/gin"

	"github.com/108356037/v1/strategy-manager/internal/kubeclient"
	"github.com/108356037/v1/strategy-manager/internal/router"
	"github.com/108356037/v1/strategy-manager/oauth"

	"github.com/108356037/v1/strategy-manager/mq"
	log "github.com/sirupsen/logrus"
)

func main() {
	// kubeclient.Init()
	// grpcclient.InitGrpcConncetion()
	// mongo.MongoConnect()
	// mq.Init()

	if os.Getenv("RUN_MODE") == "local" {
		os.Setenv("MAX_WAIT", "2")
		os.Setenv("BOOTSTRAP_SERVER", "localhost:9090")
		os.Setenv("SUBSCRIBE_TOPIC", "event")
		os.Setenv("PUBLISH_TOPIC", "event")
		os.Setenv("CONSUMER_GROUPID", "strategy-manager-svc-group")
		mq.Init()
		kubeclient.Init()
	} else {
		kubeclient.Init()
		grpcclient.InitGrpcConncetion()
		mongo.MongoConnect()
		mq.Init()
		gin.SetMode(os.Getenv("RUN_MODE"))
	}

	jobQ := make(chan string)

	go oauth.RetrieveAccessToken(jobQ)
	go func() {
		for {
			jobQ <- "tokenRefreshJob"
			time.Sleep(time.Minute * 30)
		}
	}()

	go func() {
		mq.RegisterSubscriber(os.Getenv("SUBSCRIBE_TOPIC"), mq.ResourceEventHandler)
	}()

	r := router.NewRouter()
	wt, _ := strconv.Atoi(os.Getenv("WRITE_TIMEOUT"))
	rt, _ := strconv.Atoi(os.Getenv("READ_TIMEOUT"))
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + os.Getenv("SRV_PORT"),
		WriteTimeout: time.Duration(wt) * time.Second,
		ReadTimeout:  time.Duration(rt) * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		for _, consumer := range mq.KafkaConsumerList {
			if err := consumer.Close(); err != nil {
				log.Fatal(err)
			}
			log.Infof("Closed Kafka consumer %s", consumer.Stats().ClientID)
		}

		if err := mq.KafkaProducer.Close(); err != nil {
			log.Fatal(err)
		}
		log.Infof("Closed Kafka producer %s", mq.KafkaProducer.Stats().ClientID)

		if err := grpcclient.GrpcConn.Close(); err != nil {
			log.Fatal(err)
		}
		log.Info("Closed grpc connection")

		if err := mongo.DbClient.Disconnect(*mongo.MongoCtx); err != nil {
			log.Fatal(err)
		}
		log.Info("Closed mongo connection")

		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Info(err)
	}

	log.Println("Server gracefully shutting down ...")
}
