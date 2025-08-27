package models

import (
	"time"

	"github.com/google/uuid"
)

// User model
// @Description User account information
type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
