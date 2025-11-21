package service

import (
	"context"
	"wallet-service/internal/domain"
	"wallet-service/internal/repository"

	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

type WalletService struct {
	r repository.Wallet
}

func (s *WalletService) Get(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	return s.r.Get(ctx, id)
}

func (s *WalletService) Deposit(ctx context.Context, id uuid.UUID, amount int64) (*domain.Wallet, error) {
	c, tx, err := s.r.WithTx(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			log.Error(err)
		}
	}()

	wallet, err := s.r.GetForUpdate(c, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err = wallet.Deposit(amount); err != nil {
		log.Error(err)
		return nil, err
	}

	updatedWallet, err := s.r.Update(c, wallet)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err = tx.Commit(c); err != nil {
		log.Error(err)
		return nil, err
	}

	return updatedWallet, nil
}

func (s *WalletService) Withdraw(ctx context.Context, id uuid.UUID, amount int64) (*domain.Wallet, error) {
	c, tx, err := s.r.WithTx(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			log.Error(err)
		}
	}()

	wallet, err := s.r.GetForUpdate(c, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err = wallet.Withdraw(amount); err != nil {
		log.Error(err)
		return nil, err
	}

	updatedWallet, err := s.r.Update(c, wallet)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err = tx.Commit(c); err != nil {
		log.Error(err)
		return nil, err
	}

	return updatedWallet, nil
}

func NewWalletService(r repository.Wallet) *WalletService {
	return &WalletService{
		r: r,
	}
}
