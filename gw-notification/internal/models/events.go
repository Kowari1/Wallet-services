package models

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	Deposit  EventType = "deposit"
	Withdraw EventType = "withdraw"
	Exchange EventType = "exchange"
)

type EventMessage struct {
	EventID   uuid.UUID `json:"event_id" bson:"event_id"`
	Event     EventType `json:"event" bson:"event"`
	UserID    uuid.UUID `json:"user_id" bson:"user_id"`
	WalletID  uuid.UUID `json:"wallet_id" bson:"wallet_id"`
	Amount    float64   `json:"amount" bson:"amount"`
	Currency  string    `json:"currency" bson:"currency"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Details   string    `json:"details,omitempty" bson:"details,omitempty"`
}
