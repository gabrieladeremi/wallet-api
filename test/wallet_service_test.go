package test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"wallet-api/internal/model"
	"wallet-api/internal/repo"
	"wallet-api/internal/service"
)

func newMoney(t *testing.T, cents int64) model.Money {
	m, err := model.NewMoneyFromCents(cents)
	if err != nil {
		t.Fatalf("failed to create money: %v", err)
	}
	return m
}

func TestTransferSuccess(t *testing.T) {
	ctx := context.Background()

	r := repo.NewMemoryRepo()
	s := service.NewWalletService(r)

	id1 := uuid.New().String()
	id2 := uuid.New().String()

	r.Create(ctx, &model.Wallet{ID: id1, Owner: "Alice", Balance: 10000})
	r.Create(ctx, &model.Wallet{ID: id2, Owner: "Bob", Balance: 0})

	err := s.Transfer(ctx, id1, id2, newMoney(t, 5000))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sender, _ := r.Get(ctx, id1)
	receiver, _ := r.Get(ctx, id2)

	if sender.Balance != 5000 {
		t.Errorf("wrong sender balance: expected 5000, got %d", sender.Balance)
	}

	if receiver.Balance != 5000 {
		t.Errorf("wrong receiver balance: expected 5000, got %d", receiver.Balance)
	}
}

func TestInsufficientFunds(t *testing.T) {
	ctx := context.Background()

	r := repo.NewMemoryRepo()
	s := service.NewWalletService(r)

	id1 := uuid.New().String()
	id2 := uuid.New().String()

	r.Create(ctx, &model.Wallet{ID: id1, Owner: "Alice", Balance: 1000})
	r.Create(ctx, &model.Wallet{ID: id2, Owner: "Bob", Balance: 0})

	err := s.Transfer(ctx, id1, id2, newMoney(t, 5000))
	if err == nil {
		t.Fatal("expected insufficient funds error, got nil")
	}
}

func TestInvalidAmounts(t *testing.T) {
	ctx := context.Background()
	r := repo.NewMemoryRepo()
	s := service.NewWalletService(r)

	id1 := uuid.New().String()
	id2 := uuid.New().String()
	r.Create(ctx, &model.Wallet{ID: id1, Owner: "A", Balance: 1000})
	r.Create(ctx, &model.Wallet{ID: id2, Owner: "B", Balance: 1000})

	tests := []struct {
		name   string
		amount int64
	}{
		{"zero", 0},
		{"negative", -500},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, _ := model.NewMoneyFromCents(tc.amount)
			err := s.Transfer(ctx, id1, id2, m)
			if err == nil {
				t.Fatalf("expected error for %s amount but got nil", tc.name)
			}
		})
	}
}

func TestNonexistentWallet(t *testing.T) {
	ctx := context.Background()
	r := repo.NewMemoryRepo()
	s := service.NewWalletService(r)

	err := s.Transfer(ctx, "nopeA", "nopeB", newMoney(t, 100))
	if err == nil {
		t.Fatal("expected error for nonexistent wallets, got nil")
	}
}

func TestDIUsingFailRepo(t *testing.T) {
	ctx := context.Background()

	r := repo.NewFailRepo()
	s := service.NewWalletService(r)

	err := s.Transfer(ctx, "A", "B", newMoney(t, 100))
	if err == nil {
		t.Fatal("expected error from failing repo")
	}
}

func TestConcurrentTransfers(t *testing.T) {
	ctx := context.Background()

	r := repo.NewMemoryRepo()
	s := service.NewWalletService(r)

	idA := uuid.New().String()
	idB := uuid.New().String()

	r.Create(ctx, &model.Wallet{ID: idA, Owner: "A", Balance: 100000})
	r.Create(ctx, &model.Wallet{ID: idB, Owner: "B", Balance: 0})

	transferAmount := newMoney(t, 100)

	// Run 200 concurrent transfers
	n := 200
	done := make(chan bool)

	for i := 0; i < n; i++ {
		go func() {
			_ = s.Transfer(ctx, idA, idB, transferAmount)
			done <- true
		}()
	}

	for i := 0; i < n; i++ {
		<-done
	}

	receiver, _ := r.Get(ctx, idB)

	expectedTotal := int64(n * 100)

	// receiver should get EXACTLY 200 successful transfers
	if receiver.Balance.Cents() != expectedTotal {
		t.Errorf("race condition detected: expected B=%d, got %d", expectedTotal, receiver.Balance)
	}
}
