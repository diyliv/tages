package client

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/diyliv/tages/config"
	"github.com/diyliv/tages/pkg/upload"
)

type client struct {
	logger *zap.Logger
	cfg    *config.Config
}

func NewClient(logger *zap.Logger, cfg *config.Config) *client {
	return &client{
		logger: logger,
		cfg:    cfg,
	}
}

func (c *client) Upload(path string) error {
	ctx := context.Background()

	conn, err := grpc.Dial(c.cfg.GrpcServer.Port, grpc.WithInsecure())
	if err != nil {
		c.logger.Error("Error while connecting to gRPC server: " + err.Error())
		return err
	}
	defer conn.Close()

	client := upload.NewUpload(c.logger)
	name, err := client.Upload(ctx, "test.jpeg")
	if err != nil {
		c.logger.Error("Error while calling upload method: " + err.Error())
		return err
	}
	fmt.Println(name)

	return nil
}
