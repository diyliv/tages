package main

import (
	"os"

	"github.com/diyliv/tages/config"
	"github.com/diyliv/tages/internal/server"
	"github.com/diyliv/tages/pkg/logger"
)

func main() {
	cfg := config.ReadConfig()
	logger := logger.InitLogger()
	server := server.NewServer(logger, cfg)
	if err := os.Chdir("server_saved/"); err != nil {
		panic(err)
	}
	server.StartgRPC()
}
