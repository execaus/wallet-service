package repository

import "wallet-service/internal/domain"

type Wallet interface {
	TxRepository
	Get() (*domain.Wallet, error)
	GetForUpdate() (*domain.Wallet, error)
	Update(*domain.Wallet) error
}

type Repository struct {
	Wallet
}
