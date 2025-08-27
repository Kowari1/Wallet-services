package config

import (
	"gw-exchanger/internal/pkg/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN string
	GRPCPort    int
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		logger.L.Warnw(".env not found", "err", err.Error())
	}

	port, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		port = 50051
	}

	return &Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		GRPCPort:    port,
	}, nil
}
