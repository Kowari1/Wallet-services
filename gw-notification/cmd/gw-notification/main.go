package main

import (
	"context"
	"gw-notification/internal/kafka"
	"gw-notification/internal/pkg/logger"
	"gw-notification/internal/repository"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	logger.Init()

	mongoRepo, err := repository.NewMongoRepository(
		os.Getenv("MONGO_URI"),
		os.Getenv("MONGO_DB"),
		os.Getenv("MONGO_COLLECTION"),
	)
	if err != nil {
		logger.L.Fatalw("Mongo connection failed", "error", err.Error())
	}
	defer mongoRepo.Close(context.Background())

	consumer := kafka.NewConsumer(
		os.Getenv("KAFKA_BROKER"),
		os.Getenv("KAFKA_TOPIC"),
		os.Getenv("KAFKA_GROUP_ID"),
		mongoRepo,
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			logger.L.Errorw("Consumer stopped with error", "error", err.Error())
		}
	}()

	<-ctx.Done()
	logger.L.Info("Shutting down consumer")
	_ = consumer.Close()
}
