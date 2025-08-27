package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID           uuid.UUID `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	FromCurrency string    `db:"from_currency" json:"from_currency"`
	ToCurrency   string    `db:"to_currency" json:"to_currency"`
	Amount       float64   `db:"amount" json:"amount"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
