package inventory

import (
	"encoding/json"
	"errors"
	"net/http"
	dto "purr-case/internal/dto/inventory"

	"purr-case/internal/httpapi/respond"
	inventory "purr-case/internal/service/inventory"
)

const userIdCtxKey = "userId"

type Handler struct {
	Service *inventory.Service
}

func InitHandler(service *inventory.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetUserInventory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdCtxKey).(string)
	if !ok || userID == "" {
		respond.WriteError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	items, err := h.Service.GetUserInventory(r.Context(), userID)
	if err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to get inventory")
		return
	}

	respond.WriteJSON(w, http.StatusOK, items)
}

// ConsumeInventoryItem handles user-owned item consumption.
// It validates the request body, takes userID from auth context, and delegates
// the actual quantity update to the inventory service.
func (h *Handler) ConsumeInventoryItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdCtxKey).(string)
	if !ok || userID == "" {
		respond.WriteError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	var req dto.ConsumeItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.SKU == "" {
		respond.WriteError(w, http.StatusBadRequest, "sku is required")
		return
	}
	if req.Quantity <= 0 {
		respond.WriteError(w, http.StatusBadRequest, "quantity must be a positive integer")
		return
	}

	item, err := h.Service.ConsumeItem(r.Context(), userID, req.SKU, req.Quantity)
	if err != nil {
		if errors.Is(err, inventory.ErrInsufficientInventory) {
			respond.WriteError(w, http.StatusBadRequest, "not enough items in inventory")
			return
		}
		respond.WriteError(w, http.StatusInternalServerError, "failed to consume inventory item")
		return
	}

	respond.WriteJSON(w, http.StatusOK, item)
}
