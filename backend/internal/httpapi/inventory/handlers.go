package inventory

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	inventory_dto "purr-case/internal/dto/inventory"
	items_dto "purr-case/internal/dto/items"
	"purr-case/internal/httpapi/respond"
	catalog_service "purr-case/internal/service/catalog"
	inventory_service "purr-case/internal/service/inventory"
)

const (
	userIdCtxKey = "userId"
	tokenCtxKey  = "token"
)

type Handler struct {
	Service *inventory_service.Service
	Catalog *catalog_service.Service
}

func InitHandler(service *inventory_service.Service, catalog *catalog_service.Service) *Handler {
	return &Handler{
		Service: service,
		Catalog: catalog,
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

	token, _ := r.Context().Value(tokenCtxKey).(string)
	catalogItemsBySKU := make(map[string]items_dto.Item)
	if h.Catalog != nil {
		catalogItems, err := h.Catalog.GetCatalogItems(r.Context(), token)
		if err != nil {
			log.Printf("failed to enrich inventory from catalog: %v", err)
		} else {
			catalogItemsBySKU = mapCatalogItemsBySKU(catalogItems)
		}
	}

	respond.WriteJSON(w, http.StatusOK, buildInventoryResponse(items, catalogItemsBySKU))
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

	var req inventory_dto.ConsumeItemRequest
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
		if errors.Is(err, inventory_service.ErrInsufficientInventory) {
			respond.WriteError(w, http.StatusBadRequest, "not enough items in inventory")
			return
		}
		respond.WriteError(w, http.StatusInternalServerError, "failed to consume inventory item")
		return
	}

	respond.WriteJSON(w, http.StatusOK, item)
}

var validCurrencySKUs = map[string]bool{
	"fish": true,
	"food": true,
	"yarn": true,
}

// GetCurrencyQuantity returns the quantity of a single currency item (fish/food/yarn)
func (h *Handler) GetCurrencyQuantity(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdCtxKey).(string)
	if !ok || userID == "" {
		respond.WriteError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	sku := chi.URLParam(r, "sku")
	if !validCurrencySKUs[sku] {
		respond.WriteError(w, http.StatusNotFound, "unknown currency")
		return
	}

	quantity, err := h.Service.GetItemQuantity(r.Context(), userID, sku)
	if err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to get quantity")
		return
	}

	respond.WriteJSON(w, http.StatusOK, inventory_dto.InventoryItem{
		SKU:      sku,
		Quantity: quantity,
	})
}

func mapCatalogItemsBySKU(items []items_dto.Item) map[string]items_dto.Item {
	itemsBySKU := make(map[string]items_dto.Item, len(items))
	for _, item := range items {
		if item.SKU == "" {
			continue
		}
		itemsBySKU[item.SKU] = item
	}
	return itemsBySKU
}

func buildInventoryResponse(items []inventory_dto.InventoryItem, catalogItemsBySKU map[string]items_dto.Item) inventory_dto.InventoryResponse {
	response := inventory_dto.InventoryResponse{
		Items:      make([]inventory_dto.EnrichedInventoryItem, 0, len(items)),
		Currencies: make([]inventory_dto.EnrichedInventoryItem, 0),
	}

	for _, item := range items {
		if item.Quantity <= 0 {
			continue
		}

		enriched := enrichInventoryItem(item, catalogItemsBySKU[item.SKU])
		if isCurrencyItem(enriched) {
			response.Currencies = append(response.Currencies, enriched)
			continue
		}
		response.Items = append(response.Items, enriched)
	}

	return response
}

func enrichInventoryItem(item inventory_dto.InventoryItem, catalogItem items_dto.Item) inventory_dto.EnrichedInventoryItem {
	enriched := inventory_dto.EnrichedInventoryItem{
		ID:          item.ID,
		UserID:      item.UserID,
		SKU:         item.SKU,
		Quantity:    item.Quantity,
		Name:        item.SKU,
		Description: "",
		ImageURL:    nil,
	}

	if catalogItem.SKU != "" {
		enriched.ItemID = catalogItem.ItemID
		enriched.Type = catalogItem.Type
		enriched.Name = catalogItem.Name
		enriched.Description = catalogItem.Description
		enriched.ImageURL = catalogItem.ImageURL
		enriched.Groups = catalogItem.Groups
		enriched.Content = catalogItem.Content
		enriched.CustomAttributes = catalogItem.CustomAttributes
	}

	enriched.Actions = inventoryActions(enriched)
	return enriched
}

func isCurrencyItem(item inventory_dto.EnrichedInventoryItem) bool {
	if item.Type == "virtual_currency" {
		return true
	}

	switch item.SKU {
	case "fish", "food", "yarn":
		return true
	default:
		return false
	}
}

func inventoryActions(item inventory_dto.EnrichedInventoryItem) []string {
	if strings.HasPrefix(item.SKU, "case_") {
		return []string{"open"}
	}
	return nil
}
