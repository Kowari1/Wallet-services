package storages

import (
	"context"

	"gw-currency-wallet/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

type UserStorage interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	ExistUserByUsername(ctx context.Context, username string) (bool, error)
	ExistUserByEmail(ctx context.Context, email string) (bool, error)
}

type WalletStorage interface {
	CreateWallet(ctx context.Context, wallet *models.Wallet) error
	GetWalletByUserID(ctx context.Context, userID uuid.UUID) (*models.Wallet, error)
	DepositWallet(ctx context.Context, walletID uuid.UUID, currency models.Currency, amount float64) (*models.Wallet, error)
	WithdrawWallet(ctx context.Context, walletID uuid.UUID, currency models.Currency, amount float64) (*models.Wallet, error)
}

type TransactionStorage interface {
	CreateTransaction(ctx context.Context, tx *models.Transaction) error
	ListTransactionsByUser(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error)
}
