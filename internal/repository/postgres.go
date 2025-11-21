package repository

import (
	"wallet-service/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresRepository(pool *pgxpool.Pool) (*Repository, error) {
	queries := db.New(pool)

	return &Repository{
		Wallet: NewWalletRepository(pool, queries),
	}, nil
}
