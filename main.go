package main

import (
	"log"
	"net/http"

	api "wallet-api/internal/http"
	"wallet-api/internal/repo"
	"wallet-api/internal/service"
)

func main() {
	repository := repo.NewMemoryRepo()

	svc := service.NewWalletService(repository)
	handler := api.NewHandler(svc)
	router := api.Router(handler)

	log.Println("Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
