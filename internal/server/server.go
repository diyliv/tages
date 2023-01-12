package server

import (
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/diyliv/tages/config"
	storagepb "github.com/diyliv/tages/proto/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type server struct {
	logger *zap.Logger
	cfg    *config.Config
}

func NewServer(logger *zap.Logger, cfg *config.Config) *server {
	return &server{
		logger: logger,
		cfg:    cfg,
	}
}

func (s *server) StartgRPC() {
	lis, err := net.Listen("tcp", s.cfg.GrpcServer.Port)
	if err != nil {
		s.logger.Error("Error while starting gRPC server: " + err.Error())
	}
	defer lis.Close()

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: s.cfg.GrpcServer.MaxConnectionIdle * time.Minute,
			Timeout:           s.cfg.GrpcServer.Timeout * time.Second,
			MaxConnectionAge:  s.cfg.GrpcServer.MaxConnectionAge * time.Minute,
			Time:              s.cfg.GrpcServer.Timeout * time.Minute,
		}),
	}

	server := grpc.NewServer(opts...)
	storagepb.RegisterStorageServiceServer(server, nil)

	go func() {
		if err := server.Serve(lis); err != nil {
			s.logger.Error("Error while serving: " + err.Error())
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done
	server.GracefulStop()
}
