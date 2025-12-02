package http

import (
	"github.com/go-chi/chi/v5"
)

func Router(handler *Handler) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/api/v1/health", handler.HealthCheck)
	router.Post("/api/v1/wallets", handler.CreateWallet)
	router.Get("/api/v1/wallets/{id}", handler.GetWallet)
	router.Post("/api/v1/wallets/fund", handler.FundWallet)
	router.Post("/api/v1/transfer", handler.Transfer)

	return router
}
