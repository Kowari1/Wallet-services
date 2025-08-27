package postgres

import (
	"context"
	"fmt"
	"gw-exchanger/internal/pkg/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgres(dsn string) (*PostgresDB, error) {
	if dsn == "" {
		logger.L.Warnw("POSTGRES_DSN is not set")

		return nil, fmt.Errorf("POSTGRES_DSN is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.L.Warnw("unable to create connection pool", "error", err.Error())

		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		logger.L.Warnw("database ping failed", "error", err.Error())

		return nil, err
	}

	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) Close() {
	db.Pool.Close()
}
