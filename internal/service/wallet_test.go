package service

import (
	"errors"
	"testing"
	"wallet-service/internal/domain"
	"wallet-service/internal/repository"
	mock_repository "wallet-service/internal/repository/mocks"
	"wallet-service/pkg/testdb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeposit_SuccessfulDeposit_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wallet, err := domain.NewWallet(uuid.New(), 0)
	assert.NoError(t, err)

	repo := mock_repository.NewMockWallet(ctrl)
	srv := NewWalletService(repo)

	var value int64 = 100

	mockTx := mock_repository.NewMockWalletTx(ctrl)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil).Times(1)
	mockTx.EXPECT().Rollback(gomock.Any()).AnyTimes()

	repo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)
	repo.EXPECT().GetForUpdate(t.Context(), wallet.ID()).Return(wallet, nil)
	repo.EXPECT().Update(t.Context(), wallet).Return(wallet, nil)

	finalWallet, err := srv.Deposit(t.Context(), wallet.ID(), value)
	assert.NoError(t, err)

	assert.Equal(t, value, finalWallet.Balance())
}

func TestDeposit_GetForUpdateReturnsError_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockWallet(ctrl)
	srv := NewWalletService(repo)

	walletID := uuid.New()
	expectedErr := errors.New("get for update error")

	mockTx := mock_repository.NewMockWalletTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).Times(1)

	repo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)
	repo.EXPECT().GetForUpdate(t.Context(), walletID).Return(nil, expectedErr)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

	finalWallet, err := srv.Deposit(t.Context(), walletID, 100)
	assert.Error(t, err)
	assert.Nil(t, finalWallet)
}

func TestDeposit_UpdateReturnsError_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wallet, err := domain.NewWallet(uuid.New(), 0)
	assert.NoError(t, err)

	repo := mock_repository.NewMockWallet(ctrl)
	srv := NewWalletService(repo)

	updateErr := errors.New("update balance error")

	var value int64 = 100

	mockTx := mock_repository.NewMockWalletTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).Times(1)

	repo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)
	repo.EXPECT().GetForUpdate(t.Context(), wallet.ID()).Return(wallet, nil)
	repo.EXPECT().Update(t.Context(), wallet).Return(nil, updateErr)

	finalWallet, err := srv.Deposit(t.Context(), wallet.ID(), value)
	assert.Error(t, err)
	assert.Nil(t, finalWallet)
}

func TestWithdraw_SuccessfulWithdrawal_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var initialBalance, withdrawAmount int64 = 200, 100

	wallet, err := domain.NewWallet(uuid.New(), initialBalance)
	assert.NoError(t, err)

	repo := mock_repository.NewMockWallet(ctrl)
	srv := NewWalletService(repo)

	mockTx := mock_repository.NewMockWalletTx(ctrl)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil).Times(1)
	mockTx.EXPECT().Rollback(gomock.Any()).AnyTimes()

	repo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)
	repo.EXPECT().GetForUpdate(t.Context(), wallet.ID()).Return(wallet, nil)
	repo.EXPECT().Update(t.Context(), wallet).Return(wallet, nil)

	finalWallet, err := srv.Withdraw(t.Context(), wallet.ID(), withdrawAmount)
	assert.NoError(t, err)
	assert.Equal(t, initialBalance-withdrawAmount, finalWallet.Balance())
}

func TestWithdraw_InsufficientFundsError_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	withdrawAmount := int64(100)
	initialBalance := int64(50)

	wallet, err := domain.NewWallet(uuid.New(), initialBalance)
	assert.NoError(t, err)

	repo := mock_repository.NewMockWallet(ctrl)
	srv := NewWalletService(repo)

	mockTx := mock_repository.NewMockWalletTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).Times(1)

	repo.EXPECT().WithTx(gomock.Any()).Return(t.Context(), mockTx, nil)
	repo.EXPECT().GetForUpdate(t.Context(), wallet.ID()).Return(wallet, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

	finalWallet, err := srv.Withdraw(t.Context(), wallet.ID(), withdrawAmount)
	assert.ErrorIs(t, err, domain.ErrInsufficientBalance)
	assert.Nil(t, finalWallet)
}

func TestConcurrency_TwoParallelWithdrawSecondGetsInsufficientFundsError_ReturnsError(t *testing.T) {
	t.Parallel()
	testdb.WithDB(t, func(pool *pgxpool.Pool) {
		repo, err := repository.NewPostgresRepository(pool)
		if err != nil {
			t.Fatalf("error inititalization repository: %v", err)
		}

		srv := NewWalletService(repo.Wallet)

		id, err := uuid.Parse(testdb.WalletCorrectID)
		assert.NoError(t, err)

		var amount int64 = 100
		errs := make(chan error, 2)
		done := make(chan struct{})

		for i := 0; i < 2; i++ {
			go func() {
				_, err := srv.Withdraw(t.Context(), id, amount)
				errs <- err
			}()
		}

		var isInsufficient, isSuccess bool
		for i := 0; i < 2; i++ {
			err := <-errs
			if errors.Is(err, domain.ErrInsufficientBalance) {
				isInsufficient = true
			} else if err == nil {
				isSuccess = true
			}
		}
		close(done)
		assert.True(t, isInsufficient, "one withdraw should fail with insufficient funds")
		assert.True(t, isSuccess, "one withdraw should succeed")

		w, err := repo.Get(t.Context(), id)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), w.Balance())
	})
}

func TestConcurrency_TwoParallelDepositBothSucceed_Succeed(t *testing.T) {
	t.Parallel()
	testdb.WithDB(t, func(pool *pgxpool.Pool) {
		repo, err := repository.NewPostgresRepository(pool)
		if err != nil {
			t.Fatalf("error inititalization repository: %v", err)
		}

		srv := NewWalletService(repo.Wallet)

		id, err := uuid.Parse(testdb.WalletEmptyWalletID)
		assert.NoError(t, err)

		var amount int64 = 50
		errs := make(chan error, 2)
		done := make(chan struct{})

		for i := 0; i < 2; i++ {
			go func() {
				_, err := srv.Deposit(t.Context(), id, amount)
				errs <- err
			}()
		}

		successes := 0
		for i := 0; i < 2; i++ {
			err := <-errs
			if err == nil {
				successes++
			}
		}
		close(done)
		assert.Equal(t, 2, successes, "both deposits should succeed")

		w, err := repo.Wallet.Get(t.Context(), id)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), w.Balance())
	})
}
