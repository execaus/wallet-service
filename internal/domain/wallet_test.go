package domain

import (
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

/*
Deposit(amount)

amount > 0 — увеличивает баланс

amount = 0 — ошибка

amount < 0 — ошибка

Withdraw(amount)

amount > 0 && balance >= amount — списывает

amount > balance — ошибка

amount <= 0 — ошибка
*/

func TestNewWallet_NegativeBalance_ReturnsError(t *testing.T) {
	_, err := NewWallet(uuid.New(), -1)

	assert.ErrorAs(t, err, &ErrNegativeAmount)
}

func TestDeposit_PositiveAmount_IncreasesBalance(t *testing.T) {
	var value int64 = 100
	w, _ := NewWallet(uuid.New(), 0)

	err := w.Deposit(value)

	assert.NoError(t, err)
	assert.Equal(t, value, w.Balance())
}

func TestDeposit_ZeroAmount_ReturnsError(t *testing.T) {
	var value int64
	w, _ := NewWallet(uuid.New(), 0)

	err := w.Deposit(value)

	assert.ErrorAs(t, err, &ErrZeroAmount)
}

func TestDeposit_NegativeAmount_ReturnsError(t *testing.T) {
	var value int64 = -1
	w, _ := NewWallet(uuid.New(), 0)

	err := w.Deposit(value)

	assert.ErrorAs(t, err, &ErrNegativeAmount)
}

func TestDeposit_Overflow_ReturnsError(t *testing.T) {
	var value int64 = 1
	w, _ := NewWallet(uuid.New(), math.MaxInt64)

	err := w.Deposit(value)

	assert.ErrorAs(t, err, &ErrOverflow)
}

func TestWithdraw_PositiveAmount_DecreasesBalance(t *testing.T) {
	var initValue, value, expectValue int64 = 100, 20, 80
	w, _ := NewWallet(uuid.New(), initValue)

	err := w.Withdraw(value)

	assert.NoError(t, err)
	assert.Equal(t, expectValue, w.Balance())
}

func TestWithdraw_AmountExceedsBalance_ReturnsError(t *testing.T) {
	var initValue, value int64 = 100, 101
	w, _ := NewWallet(uuid.New(), initValue)

	err := w.Withdraw(value)

	assert.ErrorAs(t, err, &ErrInsufficientBalance)
}

func TestWithdraw_ZeroAmount_ReturnsError(t *testing.T) {
	var initValue, value int64 = 100, 0
	w, _ := NewWallet(uuid.New(), initValue)

	err := w.Withdraw(value)

	assert.ErrorAs(t, err, &ErrZeroAmount)
}

func TestWithdraw_NegativeAmount_ReturnsError(t *testing.T) {
	var initValue, value int64 = 100, -1
	w, _ := NewWallet(uuid.New(), initValue)

	err := w.Withdraw(value)

	assert.ErrorAs(t, err, &ErrNegativeAmount)
}
