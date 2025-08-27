package storages

import (
	"context"
	"gw-exchanger/internal/config/models"
)

type ExchangeStorage interface {
	GetAllRates(ctx context.Context) ([]models.ExchangeRate, error)

	GetRate(ctx context.Context, fromCurrency, toCurrency string) (*models.ExchangeRate, error)
}
