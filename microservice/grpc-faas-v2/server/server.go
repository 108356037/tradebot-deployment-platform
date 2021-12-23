package server

import (
	pb "github.com/108356037/grpc/faas-deploy/v2/proto"
)

type ServerInstance struct {
	pb.UnimplementedUploadServer
	pb.UnimplementedBuildServer
	pb.UnimplementedPushServer
	pb.UnimplementedDeployServer
	pb.UnimplementedRemoveServer
	pb.UnimplementedAsyncInvokeServer
	pb.UnimplementedScheduleServer
	pb.UnimplementedDeployBotServer
	pb.UnimplementedRemoveBotServer
	pb.UnimplementedPublishBotServer
}
