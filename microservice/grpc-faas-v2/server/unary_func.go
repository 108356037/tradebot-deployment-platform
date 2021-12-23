package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/108356037/grpc/faas-deploy/v2/mq"
	pb "github.com/108356037/grpc/faas-deploy/v2/proto"
	"github.com/lithammer/shortuuid/v3"

	//"github.com/108356037/grpc/faas-deploy/v2/utils"
	log "github.com/sirupsen/logrus"
)

func (s *ServerInstance) BuildFunc(ctx context.Context, req *pb.BuildReq) (*pb.BuildRes, error) {

	// 1. rm -rf /tmp/faasCode/USER_NS/funcName && mkdir /tmp/faasCode/USER_NS/funcName
	// 2. extract tar to /tmp/faasCode/USER_NS/funcName
	// 3. build function with extracted tar and given language

	userNs := req.GetUserNS()
	userFunc := req.GetFuncName()
	lang := req.GetLang()

	imageRepo := os.Getenv("IMAGE_REPO")
	funcImage := imageRepo + "/" + userNs + "-" + userFunc

	basePath := "/tmp/faasCode"
	tarPath := path.Join(basePath, userNs, userFunc+".tar.gz")
	funcHandlerPath := path.Join(basePath, userNs) + "/" + userFunc

	//funcDirPath := path.Join(basePath, userNs)
	// tempDir := shortuuid.New()
	//runCmd(exec.Command("rm", "-rf", userFunc+".yml"))

	if err := runCmd(exec.Command("rm", "-rf", funcHandlerPath)); err != nil {
		log.Error(err.Error())
		return &pb.BuildRes{
			Message: fmt.Sprintf("error in removing dir %s: %s", funcHandlerPath, err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	if err := runCmd(exec.Command("mkdir", "-p", funcHandlerPath)); err != nil {
		log.Error(err.Error())
		return &pb.BuildRes{
			Message: fmt.Sprintf("error in mkdir %s: %s", funcHandlerPath, err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	if err := runCmd(exec.Command("tar", "-xvzf", tarPath, "-C", funcHandlerPath)); err != nil {
		log.Error(err.Error())
		return &pb.BuildRes{
			Message: fmt.Sprintf("error in decompressing tar file %s: %s", tarPath, err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	// runCmd(exec.Command("faas-cli", "new", userFunc, "--lang", lang,
	// 	"--cpu-request", "100m", "--cpu-limit", "200m",
	// 	"--memory-request", "128Mi", "--memory-limit", "256Mi", "--handler", tempDir))

	// sedQ := fmt.Sprintf("s\\handler: .*\\handler: %s\\g", funcDirPath+"/"+userFunc)
	// runCmd(exec.Command("sed", "-i", sedQ, userFunc+".yml"))

	// sedQ = fmt.Sprintf("s\\image: .*\\image: %s\\g", funcImage)
	// runCmd(exec.Command("sed", "-i", sedQ, userFunc+".yml"))

	// err := runCmd(exec.Command("faas-cli", "build", "-f", userFunc+".yml"))

	buildCmd := fmt.Sprintf("faas-cli build --image %s --name %s --lang %s --handler %s", funcImage, userFunc, lang, funcHandlerPath)
	err := runCmd(exec.Command("bash", "-c", buildCmd))
	if err != nil {
		log.Error(err.Error())
		return &pb.BuildRes{
			Message: "faas-cli build error: " + err.Error(),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	res := &pb.BuildRes{
		Message: fmt.Sprintf("Successfully builded image: %s", funcImage),
		Code:    pb.StatusCode_Ok,
	}
	// runCmd(exec.Command("rm", "-rf", tempDir))

	log.Info(res)
	return res, nil
}

func (s *ServerInstance) DeployFunc(ctx context.Context, req *pb.DeployReq) (*pb.DeployRes, error) {

	userNs := req.GetUserNS()
	userFunc := req.GetFuncName()

	imageRepo := os.Getenv("IMAGE_REPO")
	funcImage := imageRepo + "/" + userNs + "-" + userFunc
	//basePath := "/tmp/faasCode"
	//funcDirPath := path.Join(basePath, req.GetUserNS(), req.GetFuncName())

	findCmd := fmt.Sprintf("docker image ls | grep %s", funcImage)
	err := runCmd(exec.Command("bash", "-c", findCmd))
	if err != nil {
		log.Error(err.Error())
		return &pb.DeployRes{
			Message: fmt.Sprintf("Failed to deployed strategy: %s", err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	// err := runCmd(exec.Command("faas-cli", "deploy", "-n", userNs, "--env", "write_debug=true", "--yaml", userFunc+".yml"))
	deployCmd := fmt.Sprintf("faas-cli deploy --image %s --name %s -n %s --env write_debug=true --env prefix_logs=false --env write_timeout=120m --env read_timeout=120m --env exec_timeout=120m --readonly", funcImage, userFunc, userNs)
	err = runCmd(exec.Command("bash", "-c", deployCmd))
	if err != nil {
		log.Error(err.Error())
		return &pb.DeployRes{
			Message: fmt.Sprintf("Failed to deployed strategy %s at namespace %s", funcImage, userNs),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	res := &pb.DeployRes{
		Message: fmt.Sprintf("Deployed strategy %s at namespace %s", funcImage, userNs),
		Code:    pb.StatusCode_Ok,
	}
	log.Info(res)

	// This block is for publishing event
	event := mq.ResourceEvent{
		BasicEvent: mq.BasicEvent{
			EventId:    shortuuid.New(),
			OccurredAt: time.Now().Format(time.RFC3339),
		},
		ResourceEventType:  mq.ResourceCreate,
		UserId:             userNs,
		TargetResourceType: mq.Strategy,
		ResourceEventInfo: map[string]string{
			"strategy": userFunc,
		},
	}
	mqPayload, _ := json.Marshal(event)

	err = mq.PublishMsgNoKey(os.Getenv("PUBLISH_TOPIC"), mqPayload)
	if err != nil {
		mq.PublishMsgRetryQ(mqPayload)
		log.Warnf("Published event %s to retry queue", event.EventId)
	}
	log.Infof("Successfully published event %s", event.EventId)
	// This block is for publishing event

	return res, nil
}

func (s *ServerInstance) ScheduleFunc(ctx context.Context, req *pb.ScheduleReq) (*pb.ScheduleRes, error) {
	userNs := req.GetUserNS()
	userFunc := req.GetFuncName()
	imageRepo := os.Getenv("IMAGE_REPO")
	funcImage := imageRepo + "/" + userNs + "-" + userFunc
	crontabSchedule := "schedule=" + req.GetSchedule()

	findFuncCmd := fmt.Sprintf("faas-cli ls -n %s | grep %s", userNs, userFunc)
	err := runCmd(exec.Command("bash", "-c", findFuncCmd))
	if err != nil {
		log.Error(err.Error())
		return &pb.ScheduleRes{
			Message: fmt.Sprintf("Failed to schedule strategy: %s not found at ns %s", userFunc, userNs),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	deployCmd := fmt.Sprintf("faas-cli deploy --image %s --name %s -n %s --env write_debug=true --annotation topic=cron-function --annotation \"%s\" --readonly", funcImage, userFunc, userNs, crontabSchedule)
	err = runCmd(exec.Command("bash", "-c", deployCmd))
	if err != nil {
		log.Error(err.Error())
		return &pb.ScheduleRes{
			Message: fmt.Sprintf("Failed to schedule strategy %s at namespace %s: %s", funcImage, userNs, err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	res := &pb.ScheduleRes{
		Message: fmt.Sprintf("Scheduled strategy %s at namespace %s", funcImage, userNs),
		Code:    pb.StatusCode_Ok,
	}
	log.Info(res)

	// This block is for publishing event
	event := mq.ResourceEvent{
		BasicEvent: mq.BasicEvent{
			EventId: shortuuid.New(),

			OccurredAt: time.Now().Format(time.RFC3339),
		},
		ResourceEventType:  mq.ResourceUpdate,
		UserId:             userNs,
		TargetResourceType: mq.Strategy,
		ResourceEventInfo: map[string]string{
			"strategy": userFunc,
		},
		ResourceUpdateInfo: map[string]interface{}{
			"schedule": req.GetSchedule(),
		},
	}
	mqPayload, _ := json.Marshal(event)
	err = mq.PublishMsgNoKey(os.Getenv("PUBLISH_TOPIC"), mqPayload)
	if err != nil {
		mq.PublishMsgRetryQ(mqPayload)
		log.Warnf("Published event %s to retry queue", event.EventId)
	}
	log.Infof("Successfully published event %s", event.EventId)
	// This block is for publishing event

	return res, nil
}

func (s *ServerInstance) PushFunc(ctx context.Context, req *pb.PushReq) (*pb.PushRes, error) {
	userNs := req.GetUserNS()
	userFunc := req.GetFuncName()

	imageRepo := os.Getenv("IMAGE_REPO")
	funcImage := imageRepo + "/" + userNs + "-" + userFunc

	if err := runCmd(exec.Command("docker", "push", funcImage)); err != nil {
		log.Error(err.Error())
		return &pb.PushRes{
			Message: err.Error(),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	res := &pb.PushRes{
		Message: fmt.Sprintf("image %s successfully pushed", funcImage),
		Code:    pb.StatusCode_Ok,
	}
	log.Info(res)
	return res, nil
}

func (s *ServerInstance) RemoveFunc(ctx context.Context, req *pb.RemoveReq) (*pb.RemoveRes, error) {
	userNs := req.GetUserNS()
	userFunc := req.GetFuncName()

	err := runCmd(exec.Command("faas-cli", "remove", userFunc, "-n", userNs))
	if err != nil {
		log.Error(err.Error())
		return &pb.RemoveRes{
			Message: fmt.Sprintf("Failed to remove strategy %s: %s", userFunc, err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	res := &pb.RemoveRes{
		Message: fmt.Sprintf("Removed strategy %s at namespace %s", userFunc, userNs),
		Code:    pb.StatusCode_Ok,
	}
	log.Info(res)

	// This block is for publishing event
	event := mq.ResourceEvent{
		BasicEvent: mq.BasicEvent{
			EventId:    shortuuid.New(),
			OccurredAt: time.Now().Format(time.RFC3339),
		},
		ResourceEventType:  mq.ResourceDelete,
		UserId:             userNs,
		TargetResourceType: mq.Strategy,
		ResourceEventInfo: map[string]string{
			"strategy": userFunc,
		},
	}
	mqPayload, _ := json.Marshal(event)

	mq.PublishMsgNoKey(os.Getenv("PUBLISH_TOPIC"), mqPayload)
	if err != nil {
		mq.PublishMsgRetryQ(mqPayload)
		log.Warnf("Published event %s to retry queue", event.EventId)
	}
	log.Infof("Successfully published event %s", event.EventId)
	// This block is for publishing event

	return res, nil
}

func (s *ServerInstance) AsyncInvokeFunc(ctx context.Context, req *pb.AsyncInvokeReq) (*pb.AsyncInvokeRes, error) {
	userNs := req.GetUserNS()
	userFunc := req.GetFuncName()

	findFuncCmd := fmt.Sprintf("faas-cli ls -n %s | grep %s", userNs, userFunc)
	err := runCmd(exec.Command("bash", "-c", findFuncCmd))
	if err != nil {
		log.Error(err.Error())
		return &pb.AsyncInvokeRes{
			Message: fmt.Sprintf("Failed to trigger: strategy %s not found", userFunc),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	funcUrl := os.Getenv("OPENFAAS_URL") + "/async-function/" + userFunc + "." + userNs
	err = runCmd(exec.Command("curl", "-X", "POST", funcUrl))
	if err != nil {
		log.Error(err.Error())
		return &pb.AsyncInvokeRes{
			Message: fmt.Sprintf("Failed to trigger strategy %s: %s", userFunc, err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	res := &pb.AsyncInvokeRes{
		Message: fmt.Sprintf("Triggered strategy %s", userFunc),
		Code:    pb.StatusCode_Ok,
	}
	return res, nil
}

// TODO: language currently only supports python3
func (s *ServerInstance) DeployBotFunc(ctx context.Context, req *pb.DeployBotReq) (*pb.DeployBotRes, error) {

	userNs := req.GetUserNS()
	botName := req.GetBotName()
	basePath := "/tmp/faasCode"
	tarPath := path.Join(basePath, userNs, botName+".tar.gz")
	botCodePath := path.Join("/tmp/faasCode", userNs) + "/" + botName

	imageRepo := os.Getenv("IMAGE_REPO")
	botImage := imageRepo + "/" + userNs + "/" + botName + "-bot"

	if err := runCmd(exec.Command("rm", "-rf", botCodePath)); err != nil {
		log.Error(err.Error())
		return &pb.DeployBotRes{
			Message: fmt.Sprintf("error in removing dir %s: %s", botCodePath, err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	if err := runCmd(exec.Command("mkdir", "-p", botCodePath)); err != nil {
		log.Error(err.Error())
		return &pb.DeployBotRes{
			Message: fmt.Sprintf("error in mkdir %s: %s", botCodePath, err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	if err := runCmd(exec.Command("tar", "-xvzf", tarPath, "-C", botCodePath)); err != nil {
		log.Error(err.Error())
		return &pb.DeployBotRes{
			Message: fmt.Sprintf("error in decompressing tar file %s: %s", tarPath, err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	buildCmd := fmt.Sprintf("docker build -f %s -t %s %s", "/app/tradebot-dockerfiles/python3.Dockerfile", botImage, botCodePath)
	err := runCmd(exec.Command("bash", "-c", buildCmd))
	if err != nil {
		log.Error(err.Error())
		return &pb.DeployBotRes{
			Message: fmt.Sprintf("Tradebot image build error: %s", err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	pushCmd := fmt.Sprintf("docker push %s", botImage)
	err = runCmd(exec.Command("bash", "-c", pushCmd))
	if err != nil {
		log.Error(err.Error())
		return &pb.DeployBotRes{
			Message: fmt.Sprintf("Tradebot image push error: %s", err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	checkCmd := fmt.Sprintf("kubectl get deployment -n %s %s", userNs, botName)
	err = runCmd(exec.Command("bash", "-c", checkCmd))
	if err != nil {
		log.Infof("Tradebot %s in namespace %s does not exist", botName, userNs)
		deployCmd := fmt.Sprintf("kubectl create deployment %s -n %s --image %s", botName, userNs, botImage)
		patchStr := "{\"spec\": {\"strategy\":{\"rollingUpdate\":{\"maxUnavailable\":1}},\"template\": {\"metadata\": {\"labels\": {\"identity\": \"tradebot\",\"faas_function\": \"exist-to-avoid-sidecar-injection\"}},\"spec\": {\"securityContext\": {\"runAsUser\": 101,\"runAsGroup\": 101}}}}}"
		patchCmd := fmt.Sprintf("kubectl patch -n %s deployment %s --patch '%s'", userNs, botName, patchStr)

		err = runCmd(exec.Command("bash", "-c", deployCmd))
		if err != nil {
			log.Error(err.Error())
			return &pb.DeployBotRes{
				Message: fmt.Sprintf("Tradebot deployment error: %s", err.Error()),
				Code:    pb.StatusCode_Failed,
			}, nil
		}
		time.Sleep(time.Millisecond * 2000)
		err = runCmd(exec.Command("bash", "-c", patchCmd))
		if err != nil {
			log.Error(err.Error())
			return &pb.DeployBotRes{
				Message: fmt.Sprintf("Tradebot deployment error: %s", err.Error()),
				Code:    pb.StatusCode_Failed,
			}, nil
		}
	} else {
		restartCmd := fmt.Sprintf("kubectl rollout restart -n %s deployment %s", userNs, botName)
		err = runCmd(exec.Command("bash", "-c", restartCmd))
		if err != nil {
			log.Error(err.Error())
			return &pb.DeployBotRes{
				Message: fmt.Sprintf("Tradebot deployment error: %s", err.Error()),
				Code:    pb.StatusCode_Failed,
			}, nil
		}
	}

	return &pb.DeployBotRes{
		Message: fmt.Sprintf("Successfully deployed tradebot %s", botName),
		Code:    pb.StatusCode_Ok,
	}, nil

}

func (s *ServerInstance) RemoveBotFunc(ctx context.Context, req *pb.RemoveBotReq) (*pb.RemoveBotRes, error) {
	botName := req.BotName
	userNs := req.UserNS

	removeCmd := fmt.Sprintf("kubectl delete deployment -n %s %s", userNs, botName)
	err := runCmd(exec.Command("bash", "-c", removeCmd))
	if err != nil {
		log.Error(err.Error())
		return &pb.RemoveBotRes{
			Message: fmt.Sprintf("Failed deleting bot %s, error: %s", botName, err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	return &pb.RemoveBotRes{
		Message: fmt.Sprintf("Deleted bot %s", botName),
		Code:    pb.StatusCode_Ok,
	}, nil
}

// TODO: currently only supports public image
func (s *ServerInstance) PublishBotFunc(ctx context.Context, req *pb.PublishBotReq) (*pb.PublishBotRes, error) {
	botName := req.BotName
	userNs := req.UserNS
	//privacy := req.Privacy
	//privacy := "public"
	oldImage := os.Getenv("IMAGE_REPO") + "/" + userNs + "/" + botName + "-bot"
	newImage := os.Getenv("STORE_IMAGE_REPO") + "/" + userNs + "/" + botName + "-bot"

	tagCmd := fmt.Sprintf("docker tag %s %s", oldImage, newImage)
	err := runCmd(exec.Command("bash", "-c", tagCmd))
	if err != nil {
		return &pb.PublishBotRes{
			Message: fmt.Sprintf("tag error: %s", err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	pushCmd := fmt.Sprintf("docker push %s", newImage)
	err = runCmd(exec.Command("bash", "-c", pushCmd))
	if err != nil {
		return &pb.PublishBotRes{
			Message: fmt.Sprintf("push error: %s", err.Error()),
			Code:    pb.StatusCode_Failed,
		}, nil
	}

	return &pb.PublishBotRes{
		Message: fmt.Sprintf("Successfully published bot %s", botName),
		Code:    pb.StatusCode_Ok,
	}, nil

}
