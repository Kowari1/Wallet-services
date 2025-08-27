package models

import "time"

type ExchangeRate struct {
	ID           int       `db:"id" json:"-"`
	FromCurrency string    `db:"from_currency" json:"from_currency"`
	ToCurrency   string    `db:"to_currency" json:"to_currency"`
	Rate         float64   `db:"rate" json:"rate"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
