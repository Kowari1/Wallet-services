package service

import (
	"context"
	"gw-exchanger/gw-exchanger/proto/exchange"
	"gw-exchanger/internal/config/storages"
	"time"
)

type ExchangeService struct {
	exchange.UnimplementedExchangeServiceServer
	repo storages.ExchangeStorage
}

func NewExchangeServiceClient(repo storages.ExchangeStorage) *ExchangeService {
	return &ExchangeService{repo: repo}
}

func (s *ExchangeService) GetRate(ctx context.Context, req *exchange.GetRateRequest) (*exchange.GetRateResponse, error) {
	rate, err := s.repo.GetRate(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		return nil, err
	}

	return &exchange.GetRateResponse{
		FromCurrency: rate.FromCurrency,
		ToCurrency:   rate.ToCurrency,
		Rate:         float64(rate.Rate),
		UpdatedAt:    rate.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *ExchangeService) GetAllRates(ctx context.Context, req *exchange.GetAllRatesRequest) (*exchange.GetAllRatesResponse, error) {
	rates, err := s.repo.GetAllRates(ctx)
	if err != nil {
		return nil, err
	}

	var resp []*exchange.GetRateResponse
	for _, rate := range rates {
		resp = append(resp, &exchange.GetRateResponse{
			FromCurrency: rate.FromCurrency,
			ToCurrency:   rate.ToCurrency,
			Rate:         float64(rate.Rate),
			UpdatedAt:    rate.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &exchange.GetAllRatesResponse{Rates: resp}, nil
}
