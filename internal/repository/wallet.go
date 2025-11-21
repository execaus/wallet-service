package repository

import (
	"wallet-service/internal/db"
	"wallet-service/internal/domain"

	"github.com/jackc/pgx/v5"
)

type WalletRepository struct {
	TxRepositoryImpl
}

func (w *WalletRepository) Get() (*domain.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletRepository) GetForUpdate() (*domain.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (w *WalletRepository) Update(wallet *domain.Wallet) error {
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
