package repo

import (
	"context"
	"sync"

	"wallet-api/internal/model"

)

type MemoryRepo struct {
	mu      sync.Mutex
	storage map[string]*model.Wallet
}

// NewMemoryRepo creates a new in-memory wallet repository
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		storage: make(map[string]*model.Wallet),
	}
}

// Create creates a new wallet
func (repo *MemoryRepo) Create(ctx context.Context, wallet *model.Wallet) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	c := *wallet
	repo.storage[wallet.ID] = &c
	return nil
}

// Get retrieves a wallet by ID
func (repo *MemoryRepo) Get(ctx context.Context, id string) (*model.Wallet, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	wallet, ok := repo.storage[id]

	if !ok {
		return nil, ErrWalletNotFound
	}

	c := *wallet
	return &c, nil
}

// Fund adds funds to a wallet
func (repo *MemoryRepo) Fund(ctx context.Context, walletId string, amount model.Money) error {
	if amount <= 0 {
		return model.ErrInvalidAmount
	}

	repo.mu.Lock()
	defer repo.mu.Unlock()

	wallet, ok := repo.storage[walletId]
	if !ok {
		return ErrWalletNotFound
	}

	wallet.Balance += amount
	return nil
}

// Transfer transfers funds from one wallet to another
func (repo *MemoryRepo) Transfer(ctx context.Context, fromID, toID string, amount model.Money) error {
	if amount <= 0 {
		return model.ErrInvalidAmount
	}

	repo.mu.Lock()
	defer repo.mu.Unlock()

	from, ok := repo.storage[fromID]
	if !ok {
		return ErrWalletNotFound
	}

	to, ok := repo.storage[toID]
	if !ok {
		return ErrWalletNotFound
	}

	if from.Balance < amount {
		return ErrInsufficientFund
	}

	from.Balance -= amount
	to.Balance += amount

	return nil
}
