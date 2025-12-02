package repo


import (
	"context"
	"errors"

	"wallet-api/internal/model"
)

type FailRepo struct{}

func NewFailRepo() *FailRepo { return &FailRepo{} }

// Create creates a new wallet
func (repo *FailRepo) Create(ctx context.Context, wallet *model.Wallet) error {
	return errors.New("fail repo: cannot create wallet")
}

// Get retrieves a wallet by ID
func (repo *FailRepo) Get(ctx context.Context, id string) (*model.Wallet, error) {
	return nil, errors.New("fail repo: cannot get wallet")
}

// Fund adds funds to a wallet
func (repo *FailRepo) Fund(ctx context.Context, walletId string, amount model.Money) error {
	return errors.New("fail repo: cannot fund wallet")
}

// Transfer transfers funds from one wallet to another
func (repo *FailRepo) Transfer(ctx context.Context, fromID, toID string, amount model.Money) error {
	return errors.New("fail repo: transfer failed")
}
