package postgres

import (
	"context"
	"fmt"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/storages"

	"github.com/google/uuid"
)

type WalletRepo struct {
	db storages.DB
}

func NewWalletRepo(db storages.DB) storages.WalletStorage {
	return &WalletRepo{db: db}
}

func (r *WalletRepo) CreateWallet(ctx context.Context, wallet *models.Wallet) error {
	_, err := r.db.Exec(ctx, `INSERT INTO wallets (id, user_id, usd, rub, eur)
	VALUES($1, $2, $3, $4, $5)`,
		wallet.ID, wallet.UserID, wallet.USD, wallet.RUB, wallet.EUR)
	if err != nil {
		return err
	}

	return err
}

func (r *WalletRepo) GetWalletByUserID(ctx context.Context, userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	row := r.db.QueryRow(ctx, `SELECT id, user_id, usd, rub, eur, updated_at
	FROM wallets
	WHERE user_id = $1`, userID)

	err := row.Scan(&wallet.ID, &wallet.UserID, &wallet.USD, &wallet.RUB, &wallet.EUR, &wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepo) DepositWallet(ctx context.Context, walletID uuid.UUID, currency models.Currency, amount float64) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, models.ErrInvalidAmount
	}
	cents := int64(amount * float64(models.CurrencyFactor))

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM wallets WHERE id = $1 FOR UPDATE`, currency),
		walletID,
	)

	var balance int64
	if err := row.Scan(&balance); err != nil {
		return nil, err
	}

	var wallet models.Wallet
	err = tx.QueryRow(ctx,
		fmt.Sprintf(`UPDATE wallets SET %s = $1, updated_at = NOW()
		WHERE id = $2 RETURNING *`, currency),
		balance+cents, walletID,
	).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.USD,
		&wallet.RUB,
		&wallet.EUR,
		&wallet.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepo) WithdrawWallet(ctx context.Context, walletID uuid.UUID, currency models.Currency, amount float64) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, models.ErrInvalidAmount
	}
	cents := int64(amount * float64(models.CurrencyFactor))

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM wallets
		WHERE id = $1 FOR UPDATE`, currency),
		walletID,
	)

	var balance int64
	if err := row.Scan(&balance); err != nil {
		return nil, err
	}

	if balance < cents {
		return nil, models.ErrInsufficientFunds
	}

	var wallet models.Wallet
	err = tx.QueryRow(ctx,
		fmt.Sprintf(`UPDATE wallets SET %s = $1, updated_at = NOW()
		WHERE id = $2 RETURNING *`, currency),
		balance-cents, walletID,
	).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.USD,
		&wallet.RUB,
		&wallet.EUR,
		&wallet.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &wallet, nil
}
