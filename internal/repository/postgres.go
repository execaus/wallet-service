package repository

import (
	"wallet-service/internal/db"

	"github.com/jackc/pgx/v5"
)

func NewPostgresRepository(conn *pgx.Conn) (*Repository, error) {
	queries := db.New(conn)

	return &Repository{
		Wallet: NewWalletRepository(conn, queries),
	}, nil
}
