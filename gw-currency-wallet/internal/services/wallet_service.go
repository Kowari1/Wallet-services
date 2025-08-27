package services

import (
	"context"
	"fmt"
	"gw-currency-wallet/gw-exchanger/proto/exchange"
	grpcClient "gw-currency-wallet/internal/grpc"
	"gw-currency-wallet/internal/kafka"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/internal/utils"
	"time"

	"github.com/google/uuid"
)

type WalletService struct {
	walletRepo     storages.WalletStorage
	exchangeClient grpcClient.ExchangeClient
	producer       kafka.ProducerInterface
	rateCache      utils.RateCacheInterface
	sem            chan struct{}
}

func NewWalletService(walletRepo storages.WalletStorage,
	exchangeClient grpcClient.ExchangeClient,
	rateCache utils.RateCacheInterface,
	producer kafka.ProducerInterface,
	maxIn int32) *WalletService {
	return &WalletService{
		walletRepo:     walletRepo,
		exchangeClient: exchangeClient,
		producer:       producer,
		rateCache:      rateCache,
		sem:            make(chan struct{}, maxIn),
	}
}

func (s *WalletService) gate() func() {
	s.sem <- struct{}{}
	return func() { <-s.sem }
}

func (s *WalletService) CreateWallet(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	wallet := models.Wallet{
		ID:        uuid.New(),
		UserID:    userID,
		USD:       0,
		RUB:       0,
		EUR:       0,
		UpdatedAt: time.Now(),
	}

	err := s.walletRepo.CreateWallet(ctx, &wallet)
	if err != nil {
		return uuid.Nil, err
	}

	return wallet.ID, nil
}

func (s *WalletService) GetWalletByUserID(ctx context.Context, userID uuid.UUID) (*models.Wallet, error) {
	return s.walletRepo.GetWalletByUserID(ctx, userID)
}

func (s *WalletService) DepositWallet(ctx context.Context, userID uuid.UUID, currency models.Currency, amount float64) (*models.Wallet, error) {
	release := s.gate()
	defer release()

	wallet, err := s.walletRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	updateWallet, err := s.walletRepo.DepositWallet(ctx, wallet.ID, currency, amount)
	if err != nil {
		return nil, err
	}

	if amount >= models.EventAmount {
		evt := models.EventMessage{
			EventID:   uuid.New(),
			Event:     models.Deposit,
			UserID:    userID,
			WalletID:  wallet.ID,
			Amount:    amount,
			Currency:  string(currency),
			Timestamp: time.Now(),
			Details:   "",
		}
		_ = s.producer.Publish(ctx, userID.String(), evt)
	}

	return updateWallet, nil
}

func (s *WalletService) WithdrawWallet(ctx context.Context, userID uuid.UUID, currency models.Currency, amount float64) (*models.Wallet, error) {
	release := s.gate()
	defer release()

	wallet, err := s.walletRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	updateWallet, err := s.walletRepo.WithdrawWallet(ctx, wallet.ID, currency, amount)
	if err != nil {
		return nil, err
	}

	if amount >= models.EventAmount {
		evt := models.EventMessage{
			EventID:   uuid.New(),
			Event:     models.Withdraw,
			UserID:    userID,
			WalletID:  wallet.ID,
			Amount:    amount,
			Currency:  string(currency),
			Timestamp: time.Now(),
			Details:   "",
		}
		_ = s.producer.Publish(ctx, userID.String(), evt)
	}

	return updateWallet, nil
}

func (s *WalletService) GetAllRates(ctx context.Context) (map[string]float64, error) {
	if cachedRates, ok := s.rateCache.GetAllRates(); ok {
		return cachedRates, nil
	}

	rates, err := s.exchangeClient.GetAllRates(ctx)
	if err != nil {
		return nil, err
	}

	return s.createCacheRates(rates), nil
}

func (s *WalletService) ExchangeCurrency(ctx context.Context, userID uuid.UUID, from, to models.Currency, amount float64) (*models.Wallet, error) {
	release := s.gate()
	defer release()

	wallet, err := s.walletRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if _, err := wallet.Withdraw(amount, from); err != nil {
		return nil, err
	}

	rate, exists := s.rateCache.GetRate(string(from), string(to))
	if !exists {
		rates, err := s.exchangeClient.GetAllRates(ctx)
		if err != nil {
			return nil, err
		}
		s.createCacheRates(rates)

		rate, exists = s.rateCache.GetRate(string(from), string(to))
		if !exists {
			return nil, fmt.Errorf("rate form %s to %s not found even after refresh", from, to)
		}
	}

	converted := amount * rate

	if _, err = s.walletRepo.WithdrawWallet(ctx, wallet.ID, from, amount); err != nil {
		return nil, err
	}
	updatedWallet, err := s.walletRepo.DepositWallet(ctx, wallet.ID, to, converted)
	if err != nil {
		return nil, err
	}

	if amount >= models.EventAmount {
		evt := models.EventMessage{
			EventID:   uuid.New(),
			Event:     models.Exchange,
			UserID:    userID,
			WalletID:  wallet.ID,
			Amount:    amount,
			Currency:  fmt.Sprintf("%s->%s", from, to),
			Timestamp: time.Now(),
			Details:   "",
		}
		_ = s.producer.Publish(ctx, userID.String(), evt)
	}

	return updatedWallet, nil
}

func (s *WalletService) createCacheRates(rates []*exchange.GetRateResponse) map[string]float64 {
	mapRates := make(map[string]float64)

	for _, rate := range rates {
		key := fmt.Sprintf("From %s to %s", rate.FromCurrency, rate.ToCurrency)
		mapRates[key] = rate.Rate
	}

	if s.rateCache != nil {
		s.rateCache.UpdatedRates(mapRates)
	}

	return mapRates
}
