package inventory_dto

import items_dto "purr-case/internal/dto/items"

type InventoryItem struct {
	ID       int    `json:"id"`
	UserID   string `json:"user_id"`
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type EnrichedInventoryItem struct {
	ID          int                     `json:"id"`
	UserID      string                  `json:"user_id,omitempty"`
	SKU         string                  `json:"sku"`
	Quantity    int                     `json:"quantity"`
	ItemID      int                     `json:"item_id,omitempty"`
	Type        string                  `json:"type,omitempty"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ImageURL    *string                 `json:"image_url"`
	Groups      []items_dto.Group       `json:"groups,omitempty"`
	Content     []items_dto.ContentItem `json:"content,omitempty"`
	Actions     []string                `json:"actions,omitempty"`
}

type InventoryResponse struct {
	Items      []EnrichedInventoryItem `json:"items"`
	Currencies []EnrichedInventoryItem `json:"currencies"`
}

type ConsumeItemRequest struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}
