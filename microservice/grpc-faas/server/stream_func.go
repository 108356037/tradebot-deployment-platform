package server

import (
	"fmt"
	"io"
	"os"
	"path"

	pb "github.com/108356037/grpc/faas-deploy/v2/proto"
	log "github.com/sirupsen/logrus"
)

func writeToFp(fp *os.File, data []byte) error {
	w := 0
	n := len(data)
	for {
		nw, err := fp.Write(data[w:])
		if err != nil {
			return err
		}
		w += nw
		if nw >= n {
			return nil
		}
	}
}

func createFilePtr(userns, funcname string) (*os.File, error) {

	if funcname != "" {
		basePath := "/tmp/faasCode"
		if err := os.MkdirAll(basePath+"/"+userns, os.ModePerm); err != nil {
			return nil, err
		}
		filepath := path.Join(basePath, userns, funcname)
		fp, err := os.Create(filepath)
		if err != nil {
			return nil, fmt.Errorf("error in creating file ptr: %s", err.Error())
		}
		return fp, nil
	}

	return nil, fmt.Errorf("funcName must be provided")
}

func (s *ServerInstance) UploadFile(stream pb.Upload_UploadFileServer) error {

	firstChunk := true
	var fp *os.File

	for {
		data, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// streaming arrived at end of packet, break out of loop
				break
			}
			return stream.SendAndClose(
				&pb.UploadRes{
					Message: fmt.Sprintf("Error in streaming data: %s", err.Error()),
					Code:    pb.StatusCode_Failed,
				},
			)
		}

		// first chunk also receives the funcname, so server can create a file pointer to write bytes into it
		if firstChunk {
			_fp, err := createFilePtr(data.GetUserNS(), data.GetFuncName())
			if err != nil {
				log.Error(err.Error())
				return stream.SendAndClose(&pb.UploadRes{
					Message: "Unable to create file ptr: " + err.Error(),
					Code:    pb.StatusCode_Failed,
				})
			}
			fp = _fp
			firstChunk = false
		}

		defer fp.Close()

		err = writeToFp(fp, data.GetContent())
		if err != nil {
			log.Error(err.Error())
			return stream.SendAndClose(&pb.UploadRes{
				Message: "Unable to write chunk of data to filePtr:" + err.Error(),
				Code:    pb.StatusCode_Failed,
			})
		}
	}

	log.Infof("Successfully uploaded file in server, path: %s\n", fp.Name())
	return stream.SendAndClose(&pb.UploadRes{
		Message: "Upload received with success",
		Code:    pb.StatusCode_Ok,
	})

}
