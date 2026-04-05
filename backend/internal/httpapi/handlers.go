package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct{}

// Структура для позиции в корзине
type CheckoutItem struct {
	SKU 	string  	`json:"sku"`      	// Идентификатор товара
	Name 	string  	`json:"name"`      	// Название товара
	Type 	string  	`json:"type"`      	// Тип сущности: skin, case, battlepass.
	Quantity int     	`json:"quantity"`  	// Количество единиц товара
	Price  	float64 	`json:"price"`		// Цена за единицу товара
	Currency string  	`json:"currency"`  	// Валюта (например, "USD", "EUR")
}

// Структура для запроса на создание платежа
type CreateCheckoutRequest struct {
	UserID string         `json:"userId"` // Идентификатор пользователя
	Items  []CheckoutItem `json:"items"`   // Состав корзины
}

func InitHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// Точка входа для создания платежа и получения ссылки на оплату
func (h *Handler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	var req CreateCheckoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	if req.UserID == "" {
		writeError(w, http.StatusBadRequest, "userId is required")
		return
	}

	if len(req.Items) == 0 {
		writeError(w, http.StatusBadRequest, "items are required")
		return
	}

	// Проверка каждой позиции в корзине
	for i, item := range req.Items {
		if item.SKU == "" {
			writeError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].sku is required")
			return
		}
		if item.Name == "" {
			writeError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].name is required")
			return
		}
		if item.Type == "" {
			writeError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].type is required")
			return
		}
		if item.Quantity <= 0 {
			writeError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].quantity must be a positive integer")
			return
		}
		if item.Price <= 0 {
			writeError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].price must be a positive number")
			return
		}
		if item.Currency == "" {
			writeError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].currency is required")
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"orderId": "mock-order-1",
		"status": "pending",
		"itemsCount": len(req.Items),
		"checkoutUrl": "mock",
	})
}

// Резервация точки входа для обработки вебхуков от платежного провайдера
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message":"webhook received",
	})
}

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

// writeJSON sets Content-Type: application/json and encodes body to JSON.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// writeError writes a JSON object {"error": "..."} with the given status code.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// mustJSON serializes a value to JSON. Panics on error (should not happen with DTOs).
func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic("json marshal: " + err.Error())
	}
	b = append(b, '\n') // for compatibility with json.Encoder, which adds \n
	return b
}

