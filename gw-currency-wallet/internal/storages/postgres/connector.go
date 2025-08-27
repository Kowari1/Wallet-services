package postgres

import (
	"context"
	"gw-currency-wallet/internal/pkg/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgres(dsn string, maxConns, minConns int32, connLifetime time.Duration) (*PostgresDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.L.Warnw("failed to parse dsn", "err", err.Error())
		return nil, err
	}

	config.MaxConns = maxConns
	config.MinConns = minConns
	config.MaxConnLifetime = connLifetime

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return db.Pool.Exec(ctx, sql, args...)
}

func (db *PostgresDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return db.Pool.Query(ctx, sql, args...)
}

func (db *PostgresDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return db.Pool.QueryRow(ctx, sql, args...)
}

func (db *PostgresDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.Pool.Begin(ctx)
}

func (db *PostgresDB) Close() {
	db.Pool.Close()
}
