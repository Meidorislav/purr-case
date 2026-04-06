package inventory

import (
	"net/http"

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
