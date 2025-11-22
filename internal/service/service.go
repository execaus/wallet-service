package service

import (
	"context"
	"wallet-service/internal/domain"
	"wallet-service/internal/repository"

	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Wallet interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	Deposit(ctx context.Context, id uuid.UUID, amount int64) (*domain.Wallet, error)
	Withdraw(ctx context.Context, id uuid.UUID, amount int64) (*domain.Wallet, error)
}

type Service struct {
	Wallet
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Wallet: NewWalletService(repo.Wallet),
	}
}
