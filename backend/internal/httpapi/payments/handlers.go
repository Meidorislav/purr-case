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

	// Get user ID from JWT token set by the Auth middleware
	userID, ok := r.Context().Value("userId").(string)
	if !ok || userID == "" {
		respond.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	req.UserID = userID

	if len(req.Items) == 0 {
		respond.WriteError(w, http.StatusBadRequest, "items are required")
		return
	}

	// Assuming all items must have the same currency
	expectedCurrency := req.Items[0].Currency 

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
		if item.Currency != expectedCurrency {
			respond.WriteError(w, http.StatusBadRequest, "all items must have the same currency")
			return
		}
	}

	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += float64(item.Quantity) * item.Price
	}

	// Temporary mock response. In the future, the link will be provided via token.
	resp := dto.CreateCheckoutResponse{
		OrderID:     "mock-order-1",
		Status:      "pending",
		ItemsCount:  len(req.Items),
		TotalAmount: totalAmount,
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
