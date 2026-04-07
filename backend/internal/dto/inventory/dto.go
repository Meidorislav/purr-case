package inventory_dto

type InventoryItem struct {
	ID       int    `json:"id"`
	UserID   string `json:"user_id"`
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type ConsumeItemRequest struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}
