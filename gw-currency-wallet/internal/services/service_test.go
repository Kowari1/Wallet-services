package services_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"gw-currency-wallet/gw-exchanger/proto/exchange"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/pkg/logger"
	"gw-currency-wallet/internal/services"
	"gw-currency-wallet/internal/storages/postgres"

	"github.com/google/uuid"
)

func setupTestDB(t *testing.T) *postgres.PostgresDB {
	logger.Init()
	cfg := config.Load()

	db, err := postgres.NewPostgres(cfg.TestDBURL, 200, 20, time.Minute)
	if err != nil {
		t.Fatalf("не удалось подключиться к БД: %v", err)
	}

	_, err = db.Exec(context.Background(), `
        DROP TABLE IF EXISTS wallets;
        CREATE TABLE wallets (
            id UUID PRIMARY KEY,
            user_id UUID NOT NULL,
            rub NUMERIC DEFAULT 0,
            usd NUMERIC DEFAULT 0,
            eur NUMERIC DEFAULT 0,
            updated_at TIMESTAMPTZ DEFAULT now()
        );
    `)
	if err != nil {
		t.Fatalf("ошибка создания таблицы: %v", err)
	}

	return db
}

type mockExchangeClient struct{}

func (f *mockExchangeClient) GetRate(ctx context.Context, from, to string) (float64, error) {
	return 1.0, nil
}

func (f *mockExchangeClient) GetAllRates(ctx context.Context) ([]*exchange.GetRateResponse, error) {
	return []*exchange.GetRateResponse{
		{FromCurrency: "USD", ToCurrency: "EUR", Rate: 1.00},
	}, nil
}

type mockProducer struct{}

func (m *mockProducer) Publish(ctx context.Context, key string, value interface{}) error {
	return nil
}

type mockCache struct{}

func (c *mockCache) UpdatedRates(rates map[string]float64) {

}

func (c *mockCache) GetRate(from, to string) (float64, bool) {
	return 0, true
}

func (c *mockCache) GetAllRates() (map[string]float64, bool) {
	return nil, true
}

func TestWalletService_ConcurrentWithdrawals(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewWalletRepo(db)
	mockExchange := &mockExchangeClient{}
	mockProducer := &mockProducer{}
	mockCachce := &mockCache{}
	svc := services.NewWalletService(repo, mockExchange, mockCachce, mockProducer, 100)

	userID := uuid.New()
	_, err := svc.CreateWallet(context.Background(), userID)
	if err != nil {
		t.Fatalf("ошибка создания кошелька: %v", err)
	}

	_, err = svc.DepositWallet(context.Background(), userID, models.RUB, 100000.00)
	if err != nil {
		t.Fatalf("ошибка пополнения кошелька: %v", err)
	}

	const goroutines = 1000
	const withdrawAmount = 0.1

	var wg sync.WaitGroup
	wg.Add(goroutines)

	var successCount, failCount int64
	var mu sync.Mutex

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_, err := svc.WithdrawWallet(context.Background(), userID, models.RUB, withdrawAmount)
			mu.Lock()
			if err == nil {
				successCount++
			} else if err == models.ErrInsufficientFunds {
				failCount++
			} else {
				t.Errorf("неожиданная ошибка: %v", err)
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	if successCount+failCount != goroutines {
		t.Errorf("не все операции обработаны: %d успехов + %d неудач = %d, ожидалось %d",
			successCount, failCount, successCount+failCount, goroutines)
	}

	wallet, err := svc.GetWalletByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("ошибка получения кошелька: %v", err)
	}

	if wallet.RUB < 0 {
		t.Errorf("баланс не может быть отрицательным: %f", float64(wallet.RUB))
	}
}
