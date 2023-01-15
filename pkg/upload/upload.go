package upload

import (
	"context"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	storagepb "github.com/diyliv/tages/proto/storage"
)

type upload struct {
	logger   *zap.Logger
	clientpb storagepb.StorageServiceClient
}

func NewUpload(conn grpc.ClientConnInterface, logger *zap.Logger) *upload {
	return &upload{
		logger:   logger,
		clientpb: storagepb.NewStorageServiceClient(conn),
	}
}

func (u *upload) Upload(ctx context.Context, file string) (string, error) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	defer cancel()

	stream, err := u.clientpb.Upload(ctx)
	if err != nil {
		u.logger.Error("Error while calling client api: " + err.Error())
		return "", err
	}

	userFile, err := os.Open(file)
	if err != nil {
		u.logger.Error("Error while opening file: " + err.Error())
		return "", err
	}

	buf := make([]byte, 1024) // max size 1kb per stream

	for {
		num, err := userFile.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			u.logger.Error("Error while reading data: " + err.Error())
			return "", err
		}

		if err := stream.Send(&storagepb.UploadReq{Name: file, Chunk: buf[:num]}); err != nil {
			u.logger.Error("Error while sending file chunks: " + err.Error())
			return "", err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		u.logger.Error("Error while reading response from server: " + err.Error())
		return "", err
	}

	return res.GetStatus(), nil

}
