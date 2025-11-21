package service

import (
	"context"
	"wallet-service/internal/domain"
	"wallet-service/internal/repository"

	"github.com/google/uuid"
)

type WalletService struct {
	r repository.Wallet
}

func (s *WalletService) Get(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (s *WalletService) Deposit(ctx context.Context, id uuid.UUID, amount int64) (*domain.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (s *WalletService) Withdraw(ctx context.Context, id uuid.UUID, amount int64) (*domain.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func NewWalletService(r repository.Wallet) *WalletService {
	return &WalletService{
		r: r,
	}
}
