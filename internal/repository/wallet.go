package repository

import (
	"context"
	"errors"
	"wallet-service/internal/db"
	"wallet-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

var (
	ErrWalletNotFound = errors.New("wallet not found")
)

type WalletRepository struct {
	TxRepositoryImpl
}

func (r *WalletRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	q := r.getQueries(ctx)

	row, err := q.Get(ctx, UUIDToPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		log.Error(err)

		return nil, err
	}

	wallet, err := pgWalletToDomain(&row)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return wallet, nil
}

func (r *WalletRepository) GetForUpdate(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	q := r.getQueries(ctx)

	row, err := q.GetForUpdate(ctx, UUIDToPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		log.Error(err)
		return nil, err
	}

	wallet, err := pgWalletToDomain(&row)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return wallet, nil
}

func (r *WalletRepository) Update(ctx context.Context, wallet *domain.Wallet) (*domain.Wallet, error) {
	q := r.getQueries(ctx)

	row, err := q.Update(ctx, db.UpdateParams{
		ID:      UUIDToPgUUID(wallet.ID()),
		Balance: wallet.Balance(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		log.Error(err)
		return nil, err
	}

	domainWallet, err := pgWalletToDomain(&row)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return domainWallet, nil
}

func NewWalletRepository(pool *pgxpool.Pool, queries *db.Queries) *WalletRepository {
	return &WalletRepository{
		TxRepositoryImpl{
			db: pool,
			q:  queries,
		},
	}
}

func pgWalletToDomain(pgw *db.AppWallet) (*domain.Wallet, error) {
	id, err := PgUUIDToUUID(pgw.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	wallet, err := domain.NewWallet(id, pgw.Balance)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return wallet, nil
}
