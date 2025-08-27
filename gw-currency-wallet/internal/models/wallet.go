package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Currency string

// Currency represents supported currencies
// @Enum USD, RUB, EUR
const (
	RUB Currency = "RUB"
	USD Currency = "USD"
	EUR Currency = "EUR"
)

const CurrencyFactor int64 = 100

// Wallet model
// @Description User wallet with balances in different currencies
type Wallet struct {
	ID     uuid.UUID `db:"id" json:"id"`
	UserID uuid.UUID `db:"user_id" json:"user_id"`

	USD int64 `db:"usd" json:"usd"`
	RUB int64 `db:"rub" json:"rub"`
	EUR int64 `db:"eur" json:"eur"`

	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (w *Wallet) GetBalanceByCurrency(currency Currency) (float64, error) {
	switch currency {
	case RUB:
		return float64(w.RUB) / float64(CurrencyFactor), nil
	case USD:
		return float64(w.USD) / float64(CurrencyFactor), nil
	case EUR:
		return float64(w.EUR) / float64(CurrencyFactor), nil
	default:
		return 0, ErrUnsupportedCurrency
	}
}

func (w *Wallet) GetAllBalances() map[Currency]float64 {
	return map[Currency]float64{
		RUB: float64(w.RUB) / float64(CurrencyFactor),
		USD: float64(w.USD) / float64(CurrencyFactor),
		EUR: float64(w.EUR) / float64(CurrencyFactor),
	}
}

func (w *Wallet) getTarget(currency Currency) (*int64, error) {
	switch currency {
	case RUB:
		return &w.RUB, nil
	case USD:
		return &w.USD, nil
	case EUR:
		return &w.EUR, nil
	default:
		return nil, ErrUnsupportedCurrency
	}
}

func (w *Wallet) Deposit(amount float64, currency Currency) (float64, error) {
	if amount <= 0 {
		return 0, ErrInvalidAmount
	}

	cents := int64(amount * float64(CurrencyFactor))
	target, err := w.getTarget(currency)
	if err != nil {
		return 0, err
	}

	*target += cents
	return float64(*target) / float64(CurrencyFactor), nil
}

func (w *Wallet) Withdraw(amount float64, currency Currency) (float64, error) {
	if amount <= 0 {
		return 0, ErrInvalidAmount
	}

	cents := int64(amount * float64(CurrencyFactor))
	target, err := w.getTarget(currency)
	if err != nil {
		return 0, err
	}

	if *target < cents {
		return 0, ErrInsufficientFunds
	}

	*target -= cents
	return float64(*target) / float64(CurrencyFactor), nil
}

var (
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrUnsupportedCurrency = errors.New("unsupported currency")
	ErrInvalidAmount       = errors.New("invalid amount")
)
