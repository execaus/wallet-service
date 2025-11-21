package repository

import (
	"context"
	"errors"
	"wallet-service/internal/db"
	"wallet-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrWalletNotFound = errors.New("wallet not found")
)

type WalletRepository struct {
	TxRepositoryImpl
}

func (w *WalletRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletRepository) GetForUpdate(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletRepository) Update(ctx context.Context, wallet *domain.Wallet) (*domain.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func NewWalletRepository(conn *pgx.Conn, queries *db.Queries) *WalletRepository {
	return &WalletRepository{
		TxRepositoryImpl{
			db: conn,
			q:  queries,
		},
	}
}
