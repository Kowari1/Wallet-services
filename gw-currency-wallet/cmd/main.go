package main

import (
	"context"
	"os/signal"
	"syscall"

	"gw-currency-wallet/internal/config"
	grpcClient "gw-currency-wallet/internal/grpc"
	"gw-currency-wallet/internal/handlers"
	"gw-currency-wallet/internal/kafka"
	"gw-currency-wallet/internal/pkg/logger"
	"gw-currency-wallet/internal/services"
	"gw-currency-wallet/internal/storages/postgres"
	"gw-currency-wallet/internal/transport/http/middleware"
	"gw-currency-wallet/internal/utils"

	_ "gw-currency-wallet/internal/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// @title           Wallet Service API
// @description     Currency wallet management service

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	logger.Init()
	cfg := config.Load()

	db, err := postgres.NewPostgres(cfg.PostgresURL, cfg.DbMaxConns, cfg.DbMinConns, cfg.DbMaxLifetime)
	if err != nil {
		logger.L.Fatalw("failed to connect postgres", "error", err.Error())
	}
	defer db.Close()

	grpcConn, err := grpc.Dial(cfg.ExchangeGRPC,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(4<<20)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                cfg.GRPCExchangeLifetime,
			Timeout:             cfg.GRPCExchangeTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		logger.L.Fatalw("failed to connect grpc exchanger", "error", err.Error())
	}
	defer grpcConn.Close()

	producer := kafka.NewProducer(cfg.KafkaBroker, cfg.KafkaTopic, cfg.KafkaBatchSize, cfg.GRPCExchangeTimeout)
	defer producer.Close()

	walletRepo := postgres.NewWalletRepo(db)
	userRepo := postgres.NewUserRepo(db)

	cache := utils.NewRateCache(cfg.CacheRatesLifetime)

	jwtManager := services.NewJWTManager(cfg.JWTSecret)
	authService := services.NewAuthService(userRepo, walletRepo, jwtManager)
	exchangeClient := grpcClient.NewExchangeAdapter(grpcConn)
	walletService := services.NewWalletService(walletRepo, exchangeClient, cache, producer, cfg.WalletMaxInflight)

	authHandler := handlers.NewAuthHandler(authService, walletService, jwtManager)
	walletHandler := handlers.NewWalletHandler(walletService)

	r := gin.Default()

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/v1/register", authHandler.Register)
	r.POST("/api/v1/login", authHandler.Login)

	authUser := r.Group("/")
	authUser.Use(middleware.JWT(cfg.JWTSecret))
	{
		authUser.GET("/api/v1/balance", walletHandler.GetWallet)
		authUser.POST("/api/v1/wallet/deposit", walletHandler.Deposit)
		authUser.POST("/api/v1/wallet/withdraw", walletHandler.Withdraw)
		authUser.POST("/api/v1/exchange", walletHandler.Exchange)
		authUser.GET("/api/v1/exchange/rates", walletHandler.GetAllRates)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := r.Run(cfg.HTTPAddr); err != nil {
			logger.L.Fatalw("failed to run HTTP server", "error", err.Error())
		}
	}()

	<-ctx.Done()
	logger.L.Info("Shutting down wallet service")
}
