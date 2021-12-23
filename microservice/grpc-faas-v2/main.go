package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/108356037/grpc/faas-deploy/v2/mq"

	pb "github.com/108356037/grpc/faas-deploy/v2/proto"
	"github.com/108356037/grpc/faas-deploy/v2/server"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

var (
	port     = "8080"
	GrpcSrv  *grpc.Server
	Listener net.Listener
)

func gracefullStop(c chan os.Signal) {
	<-c
	log.Println("Shutdown Server ...")

	err := mq.KafkaProducer.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Closed Kafka producer")

	err = mq.KafkaConsumer.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Closed Kafka consumer")

	GrpcSrv.GracefulStop()
	log.Info("Closed grpc server")
}

func main() {

	if os.Getenv("MODE") == "local" {
		os.Setenv("MAX_WAIT", "2")
		os.Setenv("BOOTSTRAP_SERVER", "192.168.99.202:9094")
		os.Setenv("SUBSCRIBE_TOPIC", "event")
		os.Setenv("PUBLISH_TOPIC", "event")
		os.Setenv("CONSUMER_GROUPID", "grpc-faas-svc-group")
	}

	if os.Getenv("SRV_PORT") != "" {
		port = os.Getenv("SRV_PORT")
	}
	lst, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error when starting up server: %v\n", err)
	}
	Listener = lst

	grpcServer := grpc.NewServer()
	srv := &server.ServerInstance{}
	pb.RegisterUploadServer(grpcServer, srv)
	pb.RegisterBuildServer(grpcServer, srv)
	pb.RegisterPushServer(grpcServer, srv)
	pb.RegisterDeployServer(grpcServer, srv)
	pb.RegisterRemoveServer(grpcServer, srv)
	pb.RegisterAsyncInvokeServer(grpcServer, srv)
	pb.RegisterScheduleServer(grpcServer, srv)
	pb.RegisterDeployBotServer(grpcServer, srv)
	pb.RegisterRemoveBotServer(grpcServer, srv)
	pb.RegisterPublishBotServer(grpcServer, srv)

	go func() {
		if err := grpcServer.Serve(lst); err != nil {
			log.Fatal(err)
		}
	}()
	GrpcSrv = grpcServer
	log.Println("Starting grpc server at port " + port)

	mq.InitPub()
	mq.InitSub()
	go func() {
		mq.RunSubscribe(mq.SimpleLogHandler)
	}()

	// gracefully shutdown block
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	gracefullStop(quit)

}
