package grpcclient

import (
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

var (
	GrpcConn *grpc.ClientConn
)

func InitGrpcConncetion() {
	conn, err := grpc.Dial("grpc-faas-service-srv.default.svc.cluster.local:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	GrpcConn = conn
}
