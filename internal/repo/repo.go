package repo

import (
	"context"
	"errors"

	"wallet-api/internal/model"
)

var (
	ErrWalletNotFound   = errors.New("wallet not found")
	ErrInsufficientFund = errors.New("insufficient funds")
)

type WalletRepository interface {
	Create(ctx context.Context, w *model.Wallet) error
	Get(ctx context.Context, id string) (*model.Wallet, error)
	Fund(ctx context.Context, walletId string, amount model.Money) error
	Transfer(ctx context.Context, fromID, toID string, amount model.Money) error
}
