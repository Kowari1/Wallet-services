package models

// RegisterRequest represents user registration data
// @Description User registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents user login data
// @Description User login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Response represents standard API response
// @Description Standard API response format
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// WalletOperationReq represents wallet operation request
// @Description Wallet deposit/withdraw request
type WalletOperationReq struct {
	Currency string  `json:"currency" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
}

// ExchangeRequest represents currency exchange request
// @Description Currency exchange operation request
type ExchangeRequest struct {
	FromCurrency Currency `json:"from_currency" binding:"required"`
	ToCurrency   Currency `json:"to_currency" binding:"required"`
	Amount       float64  `json:"amount" binding:"required,gt=0"`
}
