package domain

import "errors"

var (
	ErrNegativeAmount      = errors.New("amount cannot be negative")
	ErrZeroAmount          = errors.New("amount cannot be zero")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrOverflow            = errors.New("balance overflow")
	ErrWalletNotFound      = errors.New("wallet not found")
)
