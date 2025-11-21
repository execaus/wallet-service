package repository

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
	"wallet-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	correctID     = "3f9a1b9e-2f64-4f42-9b4d-2d1c9a5ef901"
	nonExistentID = "3f9a1b9e-2f64-4f42-9b4d-2d1c9a5ef900"
)

func TestGet_ExistWallet_ReturnsWallet(t *testing.T) {
	withDB(t, func(r *Repository) {
		expectBalance := 100
		id, err := uuid.Parse(correctID)
		assert.NoError(t, err)

		model, err := r.Get(t.Context(), id)

		assert.NoError(t, err)
		assert.NotNil(t, model)
		assert.Equal(t, expectBalance, model.Balance())
	})
}

func TestGet_NonExistentWallet_ReturnsNilNil(t *testing.T) {
	withDB(t, func(r *Repository) {
		id, err := uuid.Parse(nonExistentID)
		assert.NoError(t, err)

		model, err := r.Get(t.Context(), id)

		assert.ErrorAs(t, err, &ErrWalletNotFound)
		assert.Nil(t, model)
	})
}

func TestUpdate_CorrectModel_ReturnsUpdatedModel(t *testing.T) {
	withDB(t, func(r *Repository) {
		var value int64 = 100
		id, err := uuid.Parse(correctID)
		assert.NoError(t, err)
		model, err := r.Get(t.Context(), id)
		assert.NoError(t, err)
		err = model.Deposit(value)
		assert.NoError(t, err)

		updatedModel, err := r.Update(t.Context(), model)

		assert.NoError(t, err)
		assert.Equal(t, model.Balance()+value, updatedModel.Balance())
	})
}

func TestUpdate_NonExistentWallet_ReturnsUpdatedModel(t *testing.T) {
	withDB(t, func(r *Repository) {
		id, err := uuid.Parse(nonExistentID)
		assert.NoError(t, err)
		model, err := domain.NewWallet(id, 0)
		assert.NoError(t, err)

		updatedModel, err := r.Update(t.Context(), model)

		assert.ErrorAs(t, err, &ErrWalletNotFound)
		assert.Nil(t, updatedModel)
	})
}

func TestGetForUpdate_CorrectWallet_LocksRow(t *testing.T) {
	withDB(t, func(r *Repository) {
		id, _ := uuid.Parse(correctID)

		// Захват блокировки
		ctx1, tx1, err := r.WithTx(t.Context())
		assert.NoError(t, err)
		defer func() { _ = tx1.Rollback(ctx1) }()

		_, err = r.GetForUpdate(ctx1, id)
		assert.NoError(t, err)

		blocked := make(chan struct{})

		// Попытка взять при блокировке
		go func() {
			ctx2, tx2, err := r.WithTx(t.Context())
			assert.NoError(t, err)
			defer func() { _ = tx2.Rollback(ctx2) }()

			_, err = r.GetForUpdate(ctx2, id)
			assert.NoError(t, err)

			close(blocked)
		}()

		// Убеждаемся, что вторая транзакция заблокирована
		select {
		case <-blocked:
			t.Fatal("second tx must be blocked but it acquired lock immediately")
		case <-time.After(time.Second):
		}

		// Освобождение блокировки
		assert.NoError(t, tx1.Commit(ctx1))

		// Теперь вторая транзакция должна выполниться
		select {
		case <-blocked:
		case <-time.After(time.Second):
			t.Fatal("second tx did not complete after lock release")
		}
	})
}

func withDB(t *testing.T, fn func(r *Repository)) {
	dbName := "app"
	dbUser := "user"
	dbPassword := "pass"

	postgresContainer, err := postgres.Run(t.Context(),
		"postgres:17",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	defer func() {
		if err = testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	host, _ := postgresContainer.Host(t.Context())
	port, _ := postgresContainer.MappedPort(t.Context(), "5432")

	dsn := fmt.Sprintf("postgres://user:pass@%s:%s/app?sslmode=disable", host, port.Port())

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			t.Fatalf("failed to close sqlDB: %v", err)
		}
	}()

	if err := goose.Up(sqlDB, "../../migrations"); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	dbConn, err := pgx.Connect(t.Context(), dsn)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	if err = dbConn.Ping(t.Context()); err != nil {
		t.Fatalf("failed to ping postgres: %v", err)
	}

	repo, err := NewPostgresRepository(dbConn)
	if err != nil {
		t.Fatalf("error inititalization repository: %v", err)
	}

	fn(repo)
}
