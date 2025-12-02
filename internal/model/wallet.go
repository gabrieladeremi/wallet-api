package model

import (
	"errors"
	"fmt"
)

type Money int64

var ErrInvalidAmount = errors.New("amount must be greater than zero")

// NewMoneyFromCents creates a Money instance from cents
func NewMoneyFromCents(cents int64) (Money, error) {
	if cents <= 0 {
		return 0, ErrInvalidAmount
	}
	return Money(cents), nil
}

// Cents returns the amount in cents
func (m Money) Cents() int64 {
	return int64(m)
}

// String returns a string representation of the money
func (m Money) String() string {
	return fmt.Sprintf("%d", m)
}

// Wallet represents a user's wallet
type Wallet struct {
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Balance Money  `json:"balance"`
}

// For creating wallets via API
type CreateWalletRequest struct {
	ID    string `json:"id"`
	Owner string `json:"owner"`
}
