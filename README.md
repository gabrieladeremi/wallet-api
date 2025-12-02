# Wallet API – Backend Engineer Assessment (Go)
This project is a small but production-grade Go service that enables secure money transfers between wallets.
It is built with clean architecture principles, strong typing for money values, and dependency injection to ensure testability and flexibility.

## Features
- Create wallets with safe, integer-based balances
- Fund wallet
- Transfer money between wallets with full validation
- Dependency injection for repository implementations
- In-memory and failing repository
- Unit tests demonstrating correctness and DI flexibility

# Architecture Overview
/main        → Application entrypoint (HTTP server)

/internal
  /model         → Core business models (Wallet, Money)
  /repo          → Repository interfaces + implementations (MemoryRepo, FailRepo)
  /service       → Business logic (WalletService)
  /http          → Handlers & routing
  
/test            → Unit tests (Transfer behaviour, DI, edge cases)

The architecture follows clean layering:

- model contains business entities
- service contains business rules
- repo abstracts away persistence
- http exposes API routes
- cmd/main wires everything together

This ensures each layer has a single responsibility.

# Dependency Injection Approach

To keep the business layer independent of storage concerns, the project uses Dependency Injection (DI) through interfaces.

### Repository Interface

Design Decisions
1. Dependency Injection Approach

To ensure clean architecture, testability, and loose coupling, the service logic does not depend on concrete implementations of data storage. Instead, I used Dependency Injection (DI) through Go interfaces.

I defined a WalletRepository interface with the required persistence operations:
```markdown
### WalletRepository Interface

```go
type WalletRepository interface {
    Get(ctx context.Context, id string) (*Wallet, error)
    Update(ctx context.Context, w *Wallet) error
    Create(ctx context.Context, w *Wallet) error
}
```
Injected Into the Service
```go
func NewWalletService(repo WalletRepository) *WalletService {
    return &WalletService{repo: repo}
}
```
Why this approach?

Swap storage easily
- MemoryRepo (tests)
- FailRepo (error simulation)
- SQLRepo or RedisRepo (future)
No infrastructure leakage  
- Service code is pure and does not know how data is stored.
Improved testability
- Tests use in-memory or mock repositories.
Clean architecture
- Business logic is decoupled from application concerns.

# Money Representation
Handling money with float64 leads to precision bugs:
```go
0.1 + 0.2 != 0.3   ❌
```
This is unacceptable for financial systems.

### Solution: Custom Money Type (integer cents)
```go
type Money struct {
    cents int64
}
```
Why integer cents?
✔ Perfect precision
✔ No rounding errors
✔ Simple arithmetic
✔ Safe validation (no negative transfers)
✔ Strong domain typing (can't mix money with integers)

```go
func NewMoneyFromCents(c int64) (Money, error) {
    if c <= 0 {
        return Money{}, errors.New("amount must be positive")
    }
    return Money{cents: c}, nil
}
```
This prevents invalid values from ever entering your service.

### Transfer Logic
The wallet service implements the required business rule:
```go
func (s *WalletService) Transfer(ctx context.Context, fromID, toID string, amount Money) error
```
The flow:
- Fetch both wallets
- Validate sender has enough balance
- Deduct amount from sender
- Add amount to receiver
- Persist both updates atomically (per repository rules)

# API Endpoints
Create Wallet

POST /wallets
```go
{
  "owner": "John Doe",
}
```
Response
```
{
    "id": "80989d2e-f566-423e-9daa-fd02e1f05306",
    "owner": "Jane Doe",
    "balance": 0
}
```
Fund Wallet
POST /wallets/fund
```go
{
    "amount": 10000,
    "wallet_id": "369264bf-6737-437b-b1be-6f6bf5498087"
}
```
Response
```
{
    "id": "369264bf-6737-437b-b1be-6f6bf5498087",
    "owner": "John Doe",
    "balance": 10000
}
```

Transfer Money

POST /transfer
```go
{
  "from": "walletA",
  "to": "walletB",
  "amount_cents": 5000
}
```
Response
```
{
  "message": "transfer successful"
}
```

# Running the Service

1. Install dependencies
   ```sh
   go mod tidy
   ```
2. Run the server
   ```sh
   go run main.go
   ```
Server starts at:
```arduino
  http://localhost:8080
```

# Running Tests
```sh
go test ./...
```

Tests include:
- Successful transfers
- Insufficient funds errors
- DI tests with FailRepo
- Concurrency safety tests

## Example: Memory Repo Storage
All data during testing or development is stored in-memory, inside:
```swift
internal/repo/memory_repo.go
```
This means:

- Data resets when the service restarts
- Fast test execution
- No external DB required
