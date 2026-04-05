package items

import (
	"encoding/json"
	"fmt"
	"net/http"
	dto "purr-case/internal/dto/items"
	"purr-case/internal/httpapi/respond"

	"github.com/go-chi/chi"
)

type Handler struct {
	ItemsURL string
}

func InitHandler(merchant_id string) *Handler {
	url := fmt.Sprintf(
		"https://store.xsolla.com/api/v2/project/%s",
		merchant_id,
	)
	return &Handler{
		ItemsURL: url,
	}
}

func (h *Handler) GetTypeItems(w http.ResponseWriter, r *http.Request, itemType string) {
	url := h.ItemsURL + "/items" + itemType
	resp, err := http.Get(url)
	if err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to fetch items")
		return
	}
	defer resp.Body.Close()

	var result dto.CatalogResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to decode response")
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
	url := h.ItemsURL + "/items/sku/" + sku
	resp, err := http.Get(url)
	if err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to fetch item")
		return
	}
	defer resp.Body.Close()

	var result dto.Item
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to decode response")
		return
	}

	respond.WriteJSON(w, http.StatusOK, result)
}
