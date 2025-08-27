package postgres

import (
	"context"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/storages"

	"github.com/google/uuid"
)

type TransactionRepo struct {
	db *PostgresDB
}

func NewTransactionRepo(db *PostgresDB) storages.TransactionStorage {
	return &TransactionRepo{db: db}
}

// CreateTransaction сохраняет транзакцию в БД.
func (r *TransactionRepo) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO transactions (id, user_id, from_currency, to_currency, amount, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		tx.ID, tx.UserID, tx.FromCurrency, tx.ToCurrency, tx.Amount, tx.CreatedAt,
	)
	return err
}

// ListTransactionsByUser возвращает все транзакции пользователя.
func (r *TransactionRepo) ListTransactionsByUser(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT id, user_id, from_currency, to_currency, amount, created_at
		FROM transactions
		WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.FromCurrency, &t.ToCurrency, &t.Amount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
