package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"log/slog"
)


type TransactionHandler struct {
	service *services.TransactionService
	logger *slog.Logger
}

func NewTransactionHandler(service *services.TransactionService, logger *slog.Logger) *TransactionHandler {
	return &TransactionHandler{service: service, logger: logger}
}

func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w,"Invalid request body", http.StatusBadRequest)
		return 
	}


	//true use lock biar thread safe
	transaction, err := h.service.Checkout(req.Items, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}


