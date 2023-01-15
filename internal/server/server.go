package server

import (
	"net"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/diyliv/tages/config"
	"github.com/diyliv/tages/internal/db"
	grpcservice "github.com/diyliv/tages/internal/delivery/grpc"
	storagepb "github.com/diyliv/tages/proto/storage"
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
	s.logger.Info("Starting gRPC server on port: " + s.cfg.GrpcServer.Port)
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

	database := db.NewDb()
	grpcservice := grpcservice.NewgRPCService(s.logger, database)
	server := grpc.NewServer(opts...)
	storagepb.RegisterStorageServiceServer(server, grpcservice)

	go func() {
		if err := server.Serve(lis); err != nil {
			s.logger.Error("Error while serving: " + err.Error())
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done
	s.logger.Info("Exiting was successful")
	server.GracefulStop()
}
