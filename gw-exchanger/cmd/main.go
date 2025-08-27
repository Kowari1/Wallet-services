package main

import (
	"gw-exchanger/internal/config"
	"gw-exchanger/internal/config/storages/postgres"
	"gw-exchanger/internal/pkg/logger"
	"gw-exchanger/server/grpc"
	"gw-exchanger/service"
	"os"
)

func main() {
	logger.Init()
	cfg, err := config.Load()
	if err != nil {
		logger.L.Warnw(".env not found use default", "err", err.Error())
	}

	db, err := postgres.NewPostgres(cfg.PostgresDSN)
	if err != nil {
		logger.L.Fatal("failed to connect to db", "error", err.Error())
	}
	defer db.Close()

	repo := postgres.NewExchangeRepo(db)
	svc := service.NewExchangeServiceClient(repo)

	port := os.Getenv("EXCHANGER_PORT")
	if port == "" {
		port = "50051"
	}

	if err := grpc.RunGRPCServer(svc, port); err != nil {
		logger.L.Fatal("server failed", "error", err.Error())
	}
}
