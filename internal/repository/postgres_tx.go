package repository

import (
	"context"
	"wallet-service/internal/db"

	"github.com/jackc/pgx/v5"
)

type TxRepository interface {
	WithTx(ctx context.Context) (context.Context, *pgx.Tx, error)
	getQueries(ctx context.Context) *db.Queries
}

type TxRepositoryImpl struct {
	db *pgx.Conn
	q  *db.Queries
}
type txKeyType struct{}

var txKey = txKeyType{}

func (r *TxRepositoryImpl) WithTx(ctx context.Context) (context.Context, *pgx.Tx, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	txQueries := r.q.WithTx(tx)

	return context.WithValue(ctx, txKey, txQueries), &tx, nil
}

func (r *TxRepositoryImpl) getQueries(ctx context.Context) *db.Queries {
	if queries := ctx.Value(txKey); queries != nil {
		return queries.(*db.Queries)
	}
	return r.q
}
