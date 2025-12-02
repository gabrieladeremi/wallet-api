package service

import (
	"context"
	"errors"

	"wallet-api/internal/model"
	"wallet-api/internal/repo"
)

type WalletService struct {
	repo repo.WalletRepository
}

var ErrInvalidAmount = errors.New("invalid amount: must be > 0")

// NewWalletService creates a new WalletService
func NewWalletService(walletRepo repo.WalletRepository) *WalletService {
	return &WalletService{repo: walletRepo}
}

// CreateWallet creates a new wallet
func (service *WalletService) CreateWallet(ctx context.Context, id, owner string) (*model.Wallet, error) {
	wallet := &model.Wallet{
		ID:      id,
		Owner:   owner,
		Balance: 0,
	}

	err := service.repo.Create(ctx, wallet)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetWallet retrieves a wallet by ID
func (service *WalletService) GetWallet(ctx context.Context, id string) (*model.Wallet, error) {
	return service.repo.Get(ctx, id)
}

// FundWallet adds funds to a wallet
func (service *WalletService) FundWallet(ctx context.Context, walletId string, amount model.Money) (*model.Wallet, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	err := service.repo.Fund(ctx, walletId, amount)
	if err != nil {
		return nil, err
	}

	return service.repo.Get(ctx, walletId)
}

// Transfer transfers funds from one wallet to another
func (service *WalletService) Transfer(ctx context.Context, fromID, toID string, amount model.Money) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	return service.repo.Transfer(ctx, fromID, toID, amount)
}
