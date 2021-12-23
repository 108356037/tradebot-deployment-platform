package server_test

import (
	"context"
	"log"
	"testing"

	pb "github.com/108356037/grpc/faas-deploy/v2/proto"
	"github.com/108356037/grpc/faas-deploy/v2/server"
	"github.com/onsi/gomega"
)

// func TestUpload(t *testing.T) {
// 	// testCases := []struct {
// 	// 	name string
// 	// 	req *pb.UploadReq
// 	// 	expectErr bool
// 	// }{
// 	// 	{
// 	// 		name: "basic upload test",
// 	// 		req: &pb.UploadReq{

// 	// 		},
// 	// 		expectErr: false,
// 	// 	},
// 	// }
// 	file, err := os.Open("../pycon.tar.gz")
// 	buf := make([]byte, 64*1024)

// 	//svc := server.ServerInstance{}
// 	client := pb.NewUploadClient(new(grpc.ClientConn))
// 	streamUploader, _ := client.UploadFile(context.Background())

// 	firstChunk := true
// 	for {
// 		n, errRead := file.Read(buf)
// 		if errRead != nil {
// 			if errRead == io.EOF {
// 				errRead = nil
// 				break
// 			}
// 			log.Fatal(err)
// 		}

// 		if firstChunk {
// 			err = streamUploader.Send(&pb.UploadReq{
// 				FuncName: "pycon",
// 				Content:  buf[:n],
// 				UserNS:   "default",
// 			})
// 			firstChunk = false
// 		} else {
// 			err = streamUploader.Send(&pb.UploadReq{
// 				Content: buf[:n],
// 			})
// 		}

// 		if err != nil {
// 			log.Println("Error when streaming data, please run again...")
// 			break
// 		}
// 	}

// 	status, _ := streamUploader.CloseAndRecv()
// 	if status.Code != pb.StatusCode_Ok {
// 		t.Fatal(status)
// 	}
// }

func TestBuildFunc(t *testing.T) {
	testCases := []struct {
		name      string
		req       *pb.BuildReq
		expectErr bool
	}{
		{
			name: "basic build test",
			req: &pb.BuildReq{
				Lang:     "python3",
				FuncName: "pycon",
				UserNS:   "default",
			},
			expectErr: false,
		},
		{
			name: "wrong lang test",
			req: &pb.BuildReq{
				Lang:     "someRandLang",
				FuncName: "pycon",
				UserNS:   "default",
			},
			expectErr: true,
		},
		{
			name: "wrong func name test",
			req: &pb.BuildReq{
				Lang:     "python3",
				FuncName: "jfiejfi",
				UserNS:   "default",
			},
			expectErr: true,
		},
		{
			name: "wrong user ns test",
			req: &pb.BuildReq{
				Lang:     "python3",
				FuncName: "pycon",
				UserNS:   "fjkeoikgoe",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			//t.Parallel()
			ctx := context.Background()

			//call
			svc := server.ServerInstance{}
			status, _ := svc.BuildFunc(ctx, testCase.req)
			if testCase.expectErr {
				if status.Code != pb.StatusCode_Failed {
					log.Fatal(status)
				}
			} else {
				if status.Code != pb.StatusCode_Ok {
					log.Fatal(status)
				}
			}
		})
	}
}

func TestPushFunc(t *testing.T) {
	testCases := []struct {
		name      string
		req       *pb.PushReq
		expectErr bool
	}{
		{
			name: "basic push test",
			req: &pb.PushReq{
				FuncName: "pycon",
				UserNS:   "default",
			},
			expectErr: false,
		},
		{
			name: "wrong func name test",
			req: &pb.PushReq{
				FuncName: "kgoekgo",
				UserNS:   "default",
			},
			expectErr: true,
		},
		{
			name: "wrong user ns test",
			req: &pb.PushReq{
				FuncName: "pycon",
				UserNS:   "kgeokgojg",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			//call
			svc := server.ServerInstance{}
			status, _ := svc.PushFunc(ctx, testCase.req)
			if testCase.expectErr {
				if status.Code != pb.StatusCode_Failed {
					log.Fatal(status)
				}
			} else {
				if status.Code != pb.StatusCode_Ok {
					log.Fatal(status)
				}
			}
		})
	}
}

func TestDeployFunc(t *testing.T) {
	testCases := []struct {
		name      string
		req       *pb.DeployReq
		expectErr bool
	}{
		{
			name: "basic deploy test",
			req: &pb.DeployReq{
				FuncName: "pycon",
				UserNS:   "default",
			},
			expectErr: false,
		},
		{
			name: "wrong user ns test",
			req: &pb.DeployReq{
				FuncName: "pycon",
				UserNS:   "gkoekgoe",
			},
			expectErr: true,
		},
		{
			name: "wrong func name test",
			req: &pb.DeployReq{
				FuncName: "keokgoeg",
				UserNS:   "default",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewGomegaWithT(t)
			ctx := context.Background()

			//call
			svc := server.ServerInstance{}
			pushStatus, _ := svc.PushFunc(ctx, &pb.PushReq{
				FuncName: testCase.req.GetFuncName(),
				UserNS:   testCase.req.GetUserNS(),
			})
			status, _ := svc.DeployFunc(ctx, testCase.req)
			if testCase.expectErr {
				if testCase.name == "wrong user ns test" {
					g.Expect(pushStatus.Code).To(gomega.BeIdenticalTo(pb.StatusCode_Failed))
					g.Expect(status.Code).To(gomega.BeIdenticalTo(pb.StatusCode_Failed))
				}

				if testCase.name == "wrong func name test" {
					g.Expect(pushStatus.Code).To(gomega.BeIdenticalTo(pb.StatusCode_Failed))
					g.Expect(status.Code).To(gomega.BeIdenticalTo(pb.StatusCode_Ok))
				}

			} else {
				g.Expect(pushStatus.Code).To(gomega.BeIdenticalTo(pb.StatusCode_Ok))
				g.Expect(status.Code).To(gomega.BeIdenticalTo(pb.StatusCode_Ok))
			}
		})
	}
}
