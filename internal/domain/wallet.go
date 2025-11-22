package domain

import (
	"math"
	"sync"

	"github.com/google/uuid"
)

var walletPool = sync.Pool{
	New: func() interface{} {
		return &Wallet{}
	},
}

type Wallet struct {
	id      uuid.UUID
	balance int64
}

func NewWallet(id uuid.UUID, balance int64) (*Wallet, error) {
	if balance < 0 {
		return nil, ErrNegativeAmount
	}
	w := walletPool.Get().(*Wallet)
	w.id = id
	w.balance = balance
	return w, nil
}

func (w *Wallet) Release() {
	w.id = uuid.Nil
	w.balance = 0
	walletPool.Put(w)
}

func (w *Wallet) ID() uuid.UUID {
	return w.id
}

func (w *Wallet) Balance() int64 {
	return w.balance
}

func (w *Wallet) Deposit(amount int64) error {
	if amount == 0 {
		return ErrZeroAmount
	}
	if amount < 0 {
		return ErrNegativeAmount
	}
	if w.balance > (math.MaxInt64 - amount) {
		return ErrOverflow
	}

	w.balance += amount

	return nil
}

func (w *Wallet) Withdraw(amount int64) error {
	if amount == 0 {
		return ErrZeroAmount
	}
	if amount < 0 {
		return ErrNegativeAmount
	}
	if w.balance < amount {
		return ErrInsufficientBalance
	}

	w.balance -= amount

	return nil
}
