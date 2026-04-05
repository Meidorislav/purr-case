package payments

import (
	"encoding/json"
	"net/http"
	dto "purr-case/internal/dto/payments"
	"purr-case/internal/httpapi/respond"
	"strconv"
)

type Handler struct{}

func InitHandler() *Handler {
	return &Handler{}
}

// Entry point for creating a payment and obtaining the payment link
func (h *Handler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCheckoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.UserID == "" {
		respond.WriteError(w, http.StatusBadRequest, "userId is required")
		return
	}

	if len(req.Items) == 0 {
		respond.WriteError(w, http.StatusBadRequest, "items are required")
		return
	}

	// Validation of each item in the cart
	for i, item := range req.Items {
		if item.SKU == "" {
			respond.WriteError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].sku is required")
			return
		}
		if item.Name == "" {
			respond.WriteError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].name is required")
			return
		}
		if item.Type == "" {
			respond.WriteError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].type is required")
			return
		}
		if item.Quantity <= 0 {
			respond.WriteError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].quantity must be a positive integer")
			return
		}
		if item.Price <= 0 {
			respond.WriteError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].price must be a positive number")
			return
		}
		if item.Currency == "" {
			respond.WriteError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].currency is required")
			return
		}
	}

	// Temporary mock response. In the future, the link will be provided via token.
	resp := dto.CreateCheckoutResponse{
		OrderID:     "mock-order-1",
		Status:      "pending",
		ItemsCount:  len(req.Items),
		CheckoutURL: "https://mock-payments.local/checkout/mock-order-1",
	}

	respond.WriteJSON(w, http.StatusOK, resp)
}

// Reserved entry point for handling webhooks from the payment provider
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	respond.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "webhook received",
	})
}
