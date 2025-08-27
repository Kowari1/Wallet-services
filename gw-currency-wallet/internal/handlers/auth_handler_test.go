package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/handlers"
	"gw-currency-wallet/internal/kafka"
	"gw-currency-wallet/internal/pkg/logger"
	"gw-currency-wallet/internal/services"
	"gw-currency-wallet/internal/storages/postgres"
	"gw-currency-wallet/internal/transport/http/middleware"
	"gw-currency-wallet/internal/utils"

	"gw-currency-wallet/gw-exchanger/proto/exchange"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type mockExchangeClient struct{}

func (f *mockExchangeClient) GetRate(ctx context.Context, from, to string) (float64, error) {
	return 0.9, nil
}

func (f *mockExchangeClient) GetAllRates(ctx context.Context) ([]*exchange.GetRateResponse, error) {
	return []*exchange.GetRateResponse{
		{FromCurrency: "USD", ToCurrency: "EUR", Rate: 0.9},
	}, nil
}

func performRequest(r http.Handler, method, path, body, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func setupTestServer(t *testing.T) (*gin.Engine, *services.JWTManager) {
	logger.Init()
	cfg := config.Load()

	db, err := postgres.NewPostgres(cfg.PostgresURL, cfg.DbMaxConns, cfg.DbMinConns, cfg.DbMaxLifetime)
	if err != nil {
		t.Fatalf("failed to connect postgres: %v", err)
	}

	cache := utils.NewRateCache(cfg.CacheRatesLifetime)

	exchangeClient := &mockExchangeClient{}

	producer := kafka.NewProducer(cfg.KafkaBroker, cfg.KafkaTopic, cfg.KafkaBatchSize, cfg.KafkaBatchTimeout)

	walletRepo := postgres.NewWalletRepo(db)
	userRepo := postgres.NewUserRepo(db)

	jwtManager := services.NewJWTManager(cfg.JWTSecret)

	authService := services.NewAuthService(userRepo, walletRepo, jwtManager)
	walletService := services.NewWalletService(walletRepo, exchangeClient, cache, producer, cfg.WalletMaxInflight)

	authHandler := handlers.NewAuthHandler(authService, walletService, jwtManager)
	walletHandler := handlers.NewWalletHandler(walletService)

	r := gin.Default()
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

	return r, jwtManager
}

func TestWalletHandlers(t *testing.T) {
	r, _ := setupTestServer(t)

	w := performRequest(r, "POST", "/api/v1/register",
		`{
			"username": "TestLogin",
			"password": "12345678",
			"email": "TestEmail@mail.ru"
		}`, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("register failed: %d %s", w.Code, w.Body.String())
	}

	w = performRequest(r, "POST", "/api/v1/login",
		`{"username": "Usertest","password": "12345678"}`, "")
	if w.Code != http.StatusOK {
		t.Fatalf("login failed: %d %s", w.Code, w.Body.String())
	}

	token := gjson.Get(w.Body.String(), "data.token").String()
	if token == "" {
		t.Fatal("empty JWT token")
	}

	w = performRequest(r, "POST", "/api/v1/wallet/deposit",
		`{"currency":"USD","amount":100}`, token)
	if w.Code != http.StatusOK {
		t.Errorf("deposit failed: %d %s", w.Code, w.Body.String())
	}

	w = performRequest(r, "POST", "/api/v1/wallet/withdraw",
		`{"amount":50, "currency":"USD"}`, token)
	if w.Code != http.StatusOK {
		t.Errorf("withdraw failed: %d %s", w.Code, w.Body.String())
	}

	w = performRequest(r, "GET", "/api/v1/balance", "", token)
	if w.Code != http.StatusOK {
		t.Errorf("balance failed: %d %s", w.Code, w.Body.String())
	}

	w = performRequest(r, "POST", "/api/v1/exchange",
		`{"from_currency":"USD","to_currency":"EUR","amount":10}`, token)
	if w.Code != http.StatusOK {
		t.Errorf("exchange failed: %d %s", w.Code, w.Body.String())
	}
}
