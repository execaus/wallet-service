package repository

import (
	"context"
	"wallet-service/internal/domain"

	"github.com/google/uuid"
)

type Wallet interface {
	TxRepository
	Get(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	GetForUpdate(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	Update(ctx context.Context, wallet *domain.Wallet) (*domain.Wallet, error)
}

type Repository struct {
	Wallet
}
