package db

import (
	"context"
	"letsgo/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg config.Config) *pgxpool.Pool {
	pool, err := pgxpool.New(
		context.Background(),
		cfg.DatabaseURL,
	)
	if err != nil {
		panic(err)
	}

	return pool
}
