package postgres

import (
	"context"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/storages"

	"github.com/google/uuid"
)

type UserRepo struct {
	db *PostgresDB
}

func NewUserRepo(db *PostgresDB) storages.UserStorage {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO users (id, username, email, password_hash)
		VALUES ($1, $2, $3, $4)`,
		user.ID, user.Username, user.Email, user.PasswordHash,
	)

	return err
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	row := r.db.Pool.QueryRow(ctx,
		`SELECT id, username, email, password_hash, created_at
		FROM users WHERE username = $1`,
		username,
	)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) ExistUserByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`,
		username,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepo) ExistUserByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`,
		email,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	row := r.db.Pool.QueryRow(ctx,
		`SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE id = $1`,
		id,
	)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
