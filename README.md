Design Decisions
1. Dependency Injection Approach

To ensure clean architecture, testability, and loose coupling, the service logic does not depend on concrete implementations of data storage. Instead, I used Dependency Injection (DI) through Go interfaces.

I defined a WalletRepository interface with the required persistence operations:

type WalletRepository interface {
  Get(ctx context.Context, id string) (*Wallet, error)
  Update(ctx context.Context, w *Wallet) error
}

The WalletService struct accepts this interface as a constructor parameter:

func NewWalletService(repo WalletRepository) *WalletService {
  return &WalletService{repo: repo}
}

This design allows the service to work with any repository implementation, such as:

An in-memory repository for testing

A failing repository to validate error paths

A future database-backed repository without changing service logic.

This architecture ensures:

Separation of concerns

Full test coverage without touching external systems

Flexible evolution of storage without refactoring business logic

Essentially, DI ensures the service is completely independent of infrastructure and focuses solely on business rules.

2. Money Representation Choices

Money handling is a common source of bugs due to floating-point precision issues.
To prevent these errors, I designed a dedicated Money type.

Why Not Use float64?

float64 produces rounding errors:

0.1 + 0.2 != 0.3

float cannot reliably represent currency

For a wallet system, this is unacceptable—money must be exact.

Chosen Approach

I used a custom Money type that stores value in the smallest currency unit (cents):

type Money struct {
  cents int64
}

This provides:

Exact arithmetic (using integers)

Simple addition/subtraction

Prevention of invalid values (negative or zero)

Strong type safety (money cannot be mixed with other int64 values)

A helper ensures safe construction:

func NewMoneyFromCents(c int64) (Money, error)

This ensures all money entering the system is valid and prevents silent overflow or misuse.

Benefits
No floating-point rounding errors

Clear money semantics

Safer API with enforced validation

Cleaner tests and business logic

