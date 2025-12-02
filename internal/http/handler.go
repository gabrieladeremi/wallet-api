package http

import (
	"context"
	"encoding/json"
	"net/http"

	"wallet-api/internal/model"
	"wallet-api/internal/service"
	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"

)

type Handler struct {
	svc *service.WalletService
}

// NewHandler creates a new HTTP handler with the given wallet service
func NewHandler(s *service.WalletService) *Handler {
	return &Handler{svc: s}
}

// Create Wallet
func (h *Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var body model.CreateWalletRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.Owner == "" {
		http.Error(w, "owner are required", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()

	wallet, err := h.svc.CreateWallet(context.Background(), id, body.Owner)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)
}

//Get Wallet
func (h *Handler) GetWallet(response http.ResponseWriter, request *http.Request) {
	 id := chi.URLParam(request, "id")

	if id == "" {
		http.Error(response, "id is required", http.StatusBadRequest)
		return
	}

	wallet, err := h.svc.GetWallet(context.Background(), id)

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(response).Encode(wallet)
}

// Fund Wallet
func (h *Handler) FundWallet(response http.ResponseWriter, request *http.Request) {
	type reqBody struct {
		Amount int64 `json:"amount"`
		WalletId string `json:"wallet_id"`
	}

	var body reqBody

	_ = json.NewDecoder(request.Body).Decode(&body)

	amount, err := model.NewMoneyFromCents(body.Amount)

	if err != nil {
		http.Error(response, "invalid amount", http.StatusBadRequest)
		return
	}

	wallet, err := h.svc.FundWallet(context.Background(), body.WalletId, amount)

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(wallet)
}

// Transfer
func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		FromID string `json:"from_id"`
		ToID   string `json:"to_id"`
		Amount int64  `json:"amount"`
	}

	var body reqBody

	_ = json.NewDecoder(r.Body).Decode(&body)

	amount, err := model.NewMoneyFromCents(body.Amount)

	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	err = h.svc.Transfer(context.Background(), body.FromID, body.ToID, amount)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "transfer successful"})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok, i am working"})
}
