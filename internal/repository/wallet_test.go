package repository

import (
	"testing"
	"time"
	"wallet-service/internal/domain"
	"wallet-service/pkg/testdb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestGet_ExistWallet_ReturnsWallet(t *testing.T) {
	testdb.WithDB(t, func(pool *pgxpool.Pool) {
		repo, err := NewPostgresRepository(pool)
		if err != nil {
			t.Fatalf("error inititalization repository: %v", err)
		}

		var expectBalance int64 = 100
		id, err := uuid.Parse(testdb.WalletCorrectID)
		assert.NoError(t, err)

		model, err := repo.Get(t.Context(), id)

		assert.NoError(t, err)
		assert.NotNil(t, model)
		assert.Equal(t, expectBalance, model.Balance())
	})
}

func TestGet_NonExistentWallet_ReturnsNilNil(t *testing.T) {
	testdb.WithDB(t, func(pool *pgxpool.Pool) {
		repo, err := NewPostgresRepository(pool)
		id, err := uuid.Parse(testdb.WalletNonExistentID)
		assert.NoError(t, err)

		model, err := repo.Get(t.Context(), id)

		assert.ErrorAs(t, err, &ErrWalletNotFound)
		assert.Nil(t, model)
	})
}

func TestUpdate_CorrectModel_ReturnsUpdatedModel(t *testing.T) {
	testdb.WithDB(t, func(pool *pgxpool.Pool) {
		repo, err := NewPostgresRepository(pool)
		var value int64 = 100
		id, err := uuid.Parse(testdb.WalletCorrectID)
		assert.NoError(t, err)
		model, err := repo.Get(t.Context(), id)
		assert.NoError(t, err)
		err = model.Deposit(value)
		assert.NoError(t, err)

		updatedModel, err := repo.Update(t.Context(), model)

		assert.NoError(t, err)
		assert.Equal(t, model.Balance(), updatedModel.Balance())
	})
}

func TestUpdate_NonExistentWallet_ReturnsUpdatedModel(t *testing.T) {
	testdb.WithDB(t, func(pool *pgxpool.Pool) {
		repo, err := NewPostgresRepository(pool)
		id, err := uuid.Parse(testdb.WalletNonExistentID)
		assert.NoError(t, err)
		model, err := domain.NewWallet(id, 0)
		assert.NoError(t, err)

		updatedModel, err := repo.Update(t.Context(), model)

		assert.ErrorAs(t, err, &ErrWalletNotFound)
		assert.Nil(t, updatedModel)
	})
}

func TestGetForUpdate_CorrectWallet_LocksRow(t *testing.T) {
	testdb.WithDB(t, func(pool *pgxpool.Pool) {
		repo, err := NewPostgresRepository(pool)
		id, _ := uuid.Parse(testdb.WalletCorrectID)

		// Захват блокировки
		ctx1, tx1, err := repo.WithTx(t.Context())
		assert.NoError(t, err)
		defer func() { _ = tx1.Rollback(ctx1) }()

		_, err = repo.GetForUpdate(ctx1, id)
		assert.NoError(t, err)

		blocked := make(chan struct{})

		// Попытка взять при блокировке
		go func() {
			ctx2, tx2, err := repo.WithTx(t.Context())
			assert.NoError(t, err)
			defer func() { _ = tx2.Rollback(ctx2) }()

			_, err = repo.GetForUpdate(ctx2, id)
			assert.NoError(t, err)

			close(blocked)
		}()

		// Убеждаемся, что вторая транзакция заблокирована
		select {
		case <-blocked:
			t.Fatal("вторая транзакция должна быть заблокирована, но она сразу получила блокировку")
		case <-time.After(time.Second):
		}

		// Освобождение блокировки
		assert.NoError(t, tx1.Commit(ctx1))

		// Теперь вторая транзакция должна выполниться
		select {
		case <-blocked:
		case <-time.After(time.Second):
			t.Fatal("вторая транзакция не завершилась после освобождения блокировки")
		}
	})
}

func TestGetForUpdate_ConcurrentDeposits_CorrectBalance(t *testing.T) {
	testdb.WithDB(t, func(pool *pgxpool.Pool) {
		repo, err := NewPostgresRepository(pool)
		var deposit1, deposit2 int64 = 50, 30

		id, _ := uuid.Parse(testdb.WalletEmptyWalletID)

		// Захват блокировки
		ctx1, tx1, err := repo.WithTx(t.Context())
		assert.NoError(t, err)
		defer func() { _ = tx1.Rollback(ctx1) }()

		w1, err := repo.GetForUpdate(ctx1, id)
		assert.NoError(t, err)

		blocked := make(chan struct{})

		// Попытка взять при блокировке
		go func() {
			ctx2, tx2, err := repo.WithTx(t.Context())
			assert.NoError(t, err)
			defer func() { _ = tx2.Rollback(ctx2) }()

			w2, err := repo.GetForUpdate(ctx2, id)
			assert.NoError(t, err)

			// Пополнение первой транзакции
			assert.NoError(t, w2.Deposit(deposit2))
			_, err = repo.Update(ctx2, w2)
			assert.NoError(t, err)

			// Освобождение блокировки
			assert.NoError(t, tx2.Commit(ctx2))

			close(blocked)
		}()

		// Убеждаемся, что вторая транзакция заблокирована
		select {
		case <-blocked:
			t.Fatal("вторая транзакция должна быть заблокирована, но она сразу получила блокировку")
		case <-time.After(time.Second):
		}

		// Пополнение первой транзакции
		assert.NoError(t, w1.Deposit(deposit1))
		_, err = repo.Update(ctx1, w1)
		assert.NoError(t, err)

		// Освобождение блокировки
		assert.NoError(t, tx1.Commit(ctx1))

		// Теперь вторая транзакция должна выполниться
		select {
		case <-blocked:
		case <-time.After(time.Second):
			t.Fatal("вторая транзакция не завершилась после освобождения блокировки")
		}

		w, err := repo.Wallet.Get(t.Context(), id)
		assert.NoError(t, err)
		assert.Equal(t, deposit1+deposit2, w.Balance())
	})
}
