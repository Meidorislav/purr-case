package items

import (
	"net/http"

	"purr-case/internal/httpapi/respond"
	catalog_service "purr-case/internal/service/catalog"

	"github.com/go-chi/chi"
)

const tokenCtxKey = "token"

type Handler struct {
	Catalog *catalog_service.Service
}

func InitHandler(catalog *catalog_service.Service) *Handler {
	return &Handler{
		Catalog: catalog,
	}
}

func (h *Handler) GetTypeItems(w http.ResponseWriter, r *http.Request, itemType string) {
	token, _ := r.Context().Value(tokenCtxKey).(string)
	result, err := h.Catalog.FetchItems(r.Context(), token, itemType, r.URL.RawQuery)
	if err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to fetch items")
		return
	}

	respond.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	h.GetTypeItems(w, r, "")
}

func (h *Handler) GetVirtualItems(w http.ResponseWriter, r *http.Request) {
	h.GetTypeItems(w, r, "/virtual_items")
}

func (h *Handler) GetItemBySku(w http.ResponseWriter, r *http.Request) {
	sku := chi.URLParam(r, "sku")
	token, _ := r.Context().Value(tokenCtxKey).(string)
	result, err := h.Catalog.FetchItemBySKU(r.Context(), token, sku, r.URL.RawQuery)
	if err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to fetch item")
		return
	}

	respond.WriteJSON(w, http.StatusOK, result)
}
