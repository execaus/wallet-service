package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Wallet struct {
	id      uuid.UUID
	balance int64
}

var (
	ErrNegativeAmount      = errors.New("amount cannot be negative")
	ErrZeroAmount          = errors.New("amount cannot be zero")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrOverflow            = errors.New("balance overflow")
)

func NewWallet(id uuid.UUID, balance int64) (*Wallet, error) {
	return &Wallet{
		id:      id,
		balance: balance,
	}, nil
}

func (w *Wallet) ID() uuid.UUID {
	return w.id
}

func (w *Wallet) Balance() int64 {
	return w.balance
}

func (w *Wallet) Deposit(amount int64) error {
	// TODO
}

func (w *Wallet) Withdraw(amount int64) error {
	// TODO
}
