package config

import (
	"gw-currency-wallet/internal/pkg/logger"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr  string
	JWTSecret string

	PostgresURL   string
	TestDBURL     string
	DbMaxConns    int32
	DbMinConns    int32
	DbMaxLifetime time.Duration

	WalletMaxInflight int32

	KafkaBroker       string
	KafkaTopic        string
	KafkaBatchSize    int
	KafkaBatchTimeout time.Duration

	CacheRatesLifetime time.Duration

	ExchangeGRPC         string
	GRPCExchangeLifetime time.Duration
	GRPCExchangeTimeout  time.Duration
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		logger.L.Warn(".env not found")
	}

	cfg := &Config{
		HTTPAddr:  getEnvStr("HTTP_ADDR", ":8080"),
		JWTSecret: getEnvStr("JWT_SECRET", "supersecret"),

		PostgresURL: getEnvStr("POSTGRES_URL", "postgres://wallet_user:wallet_pass@localhost:5432/wallet_db?sslmode=disable"),

		TestDBURL: getEnvStr("TEST_DB_DSN", "postgres://test_user:test_pass@localhost:5433/gw_wallet_test?sslmode=disable"),

		DbMaxConns:    getEnvInt32("DB_MAX_CONNS", 200),
		DbMinConns:    getEnvInt32("DB_MIN_CONNS", 20),
		DbMaxLifetime: getEnvDuration("DB_MAX_LIFETIME", 5*time.Minute),

		WalletMaxInflight: getEnvInt32("WALLET_MAX_INFLIGHT", 150),

		CacheRatesLifetime: getEnvDuration("CACHE_RATES_LIFETIME", 1*time.Minute),

		KafkaBroker:       getEnvStr("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:        getEnvStr("KAFKA_TOPIC", "wallet-events"),
		KafkaBatchSize:    getEnvInt("KAFKA_BROKER", 100),
		KafkaBatchTimeout: getEnvDuration("KAFKA_BATCH_TIMEOUT", 50*time.Millisecond),

		ExchangeGRPC: getEnvStr("EXCHANGE_GRPC", "localhost:50051"),

		GRPCExchangeLifetime: getEnvDuration("GRPC_EXCHANGE_ALIVE_TIME", 30*time.Second),
		GRPCExchangeTimeout:  getEnvDuration("GRPC_EXCHANGE_ALIVE_TIMEOUT", 5*time.Second),
	}

	return cfg
}

func getEnvStr(key, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	logger.L.Warn("use default config value", def)
	return def
}

func getEnvInt(key string, def int) int {
	if s := os.Getenv(key); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			return int(v)
		}
	}
	return def
}

func getEnvInt32(key string, def int32) int32 {
	if s := os.Getenv(key); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			return int32(v)
		}
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	if s := os.Getenv(key); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			return d
		}
	}
	return def
}
