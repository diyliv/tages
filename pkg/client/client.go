package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	ratelimiter "github.com/tommy-sho/rate-limiter-grpc-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/diyliv/tages/config"
	"github.com/diyliv/tages/pkg/logger"
	"github.com/diyliv/tages/pkg/upload"
	storagepb "github.com/diyliv/tages/proto/storage"
)

func main() {
	ctx := context.Background()
	logger := logger.InitLogger()
	cfg := config.ReadConfig()

	conn, err := grpc.Dial(cfg.GrpcServer.Port,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(ratelimiter.StreamClientInterceptor(ratelimiter.NewLimiter(10))))
	if err != nil {
		logger.Error("Error while connecting to gRPC server: " + err.Error())
	}
	defer conn.Close()

	client := upload.NewUpload(conn, logger)

	files, err := ioutil.ReadDir("client_files/")
	if err != nil {
		logger.Error("Error while reading directory: " + err.Error())
	}

	for _, file := range files {
		logger.Info(fmt.Sprintf("Starting uploading: %s", file.Name()))
		name, err := client.Upload(ctx, "client_files/"+file.Name())
		if err != nil {
			logger.Error("Error while calling upload method: " + err.Error())
		}
		fmt.Println("[RESULT]", name)
	}

	logger.Info("Starting downloading files from server")
	grpcclient := storagepb.NewStorageServiceClient(conn)
	stream, err := grpcclient.SendFiles(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Error("Error while calling sendfiles() method: " + err.Error())
	}

	if err := os.Chdir("client_requested_files/"); err != nil {
		logger.Error("Error while changing directories: " + err.Error())
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Error("Error while receiving data: " + err.Error())
		}
		logger.Info("Saving file: " + resp.GetName())
		if err := ioutil.WriteFile(resp.GetName(), resp.GetChunk(), 0644); err != nil {
			logger.Error("Error while saving file: " + err.Error())
		}
	}
	logger.Info("Saving RPC was successful")

	logger.Info("Starting reading files")
	readConn, err := grpc.Dial(cfg.GrpcServer.Port, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(ratelimiter.UnaryClientInterceptor(ratelimiter.NewLimiter(100))))
	if err != nil {
		logger.Error("error while connecting: " + err.Error())
	}

	readClient := storagepb.NewStorageServiceClient(readConn)
	defer func() {
		conn.Close()
		readConn.Close()
	}()

	res, err := readClient.GetAllFiles(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Error("error while calling GetAllFiles RPC: " + err.Error())
	}

	logger.Info(fmt.Sprintf("%v\n", res.GetFile()))
}
