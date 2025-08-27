package postgres

import (
	"context"
	"gw-exchanger/internal/config/models"
	"gw-exchanger/internal/pkg/logger"

	"github.com/jackc/pgx/v5"
)

type ExchangeRepo struct {
	db *PostgresDB
}

func NewExchangeRepo(db *PostgresDB) *ExchangeRepo {
	return &ExchangeRepo{db: db}
}

func (r *ExchangeRepo) GetAllRates(ctx context.Context) ([]models.ExchangeRate, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT from_currency, to_currency, rate, updated_at
		FROM exchange_rates
		ORDER BY from_currency, to_currency`)
	if err != nil {
		logger.L.Warnw("failed to query exchange rates", "err", err.Error())

		return nil, err
	}
	defer rows.Close()

	var rates []models.ExchangeRate
	for rows.Next() {
		var rate models.ExchangeRate
		if err := rows.Scan(
			&rate.FromCurrency,
			&rate.ToCurrency,
			&rate.Rate,
			&rate.UpdatedAt,
		); err != nil {
			logger.L.Warnw("failed to scan rate row", "err", err.Error())

			return nil, err
		}

		rates = append(rates, rate)
	}

	return rates, nil
}

func (r *ExchangeRepo) GetRate(ctx context.Context, fromCurrency, toCurrency string) (*models.ExchangeRate, error) {
	var rate models.ExchangeRate

	err := r.db.Pool.QueryRow(ctx,
		`SELECT from_currency, to_currency, rate, updated_at
		FROM exchange_rates
		WHERE (from_currency = $1 AND to_currency = $2)
		OR (from_currency = $2 AND to_currency = $1)`, fromCurrency, toCurrency).Scan(
		&rate.FromCurrency,
		&rate.ToCurrency,
		&rate.Rate,
		&rate.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.L.Warnw("failed to get exchange rate", "err", err.Error())

		return nil, err
	}

	if rate.FromCurrency != fromCurrency {
		rate.Rate = 1 / rate.Rate
	}

	return &rate, nil
}

func ToFloatRate(rate int64) float64 {
	return float64(rate) / float64(100)
}
