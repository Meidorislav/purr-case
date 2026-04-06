package items_dto

import "encoding/json"

type Price struct {
	Amount                string `json:"amount"`
	AmountWithoutDiscount string `json:"amount_without_discount"`
	Currency              string `json:"currency"`
}

type CalculatedPrice struct {
	Amount                string `json:"amount"`
	AmountWithoutDiscount string `json:"amount_without_discount"`
}

type VirtualPrice struct {
	ItemID                int              `json:"item_id"`
	SKU                   string           `json:"sku"`
	Name                  string           `json:"name"`
	Type                  string           `json:"type"`
	Description           string           `json:"description"`
	ImageURL              *string          `json:"image_url"`
	Amount                int              `json:"amount"`
	AmountWithoutDiscount int              `json:"amount_without_discount"`
	CalculatedPrice       *CalculatedPrice `json:"calculated_price,omitempty"`
	IsDefault             bool             `json:"is_default"`
}

type Group struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}

type Period struct {
	DateFrom  string `json:"date_from"`
	DateUntil string `json:"date_until"`
}

type Consumable struct {
	UsagesCount int `json:"usages_count"`
}

type InventoryOptions struct {
	Consumable       *Consumable     `json:"consumable"`
	ExpirationPeriod json.RawMessage `json:"expiration_period"`
}

type ContentItem struct {
	ItemID           int               `json:"item_id"`
	SKU              string            `json:"sku"`
	Type             string            `json:"type"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	ImageURL         *string           `json:"image_url"`
	Price            *Price            `json:"price"`
	VirtualPrices    []VirtualPrice    `json:"virtual_prices,omitempty"`
	Attributes       []interface{}     `json:"attributes,omitempty"`
	IsFree           bool              `json:"is_free"`
	Groups           []Group           `json:"groups,omitempty"`
	VirtualItemType  string            `json:"virtual_item_type,omitempty"`
	InventoryOptions *InventoryOptions `json:"inventory_options,omitempty"`
	Quantity         int               `json:"quantity,omitempty"`
	Limits           json.RawMessage   `json:"limits"`
}

type Item struct {
	ItemID            int               `json:"item_id"`
	SKU               string            `json:"sku"`
	Type              string            `json:"type"`
	Name              string            `json:"name"`
	BundleType        string            `json:"bundle_type,omitempty"`
	Description       string            `json:"description"`
	ImageURL          *string           `json:"image_url"`
	Price             *Price            `json:"price"`
	VirtualPrices     []VirtualPrice    `json:"virtual_prices"`
	CanBeBought       bool              `json:"can_be_bought"`
	Promotions        []interface{}     `json:"promotions,omitempty"`
	Limits            json.RawMessage   `json:"limits,omitempty"`
	Periods           []Period          `json:"periods"`
	Attributes        []interface{}     `json:"attributes,omitempty"`
	IsFree            bool              `json:"is_free"`
	Groups            []Group           `json:"groups"`
	VirtualItemType   string            `json:"virtual_item_type,omitempty"`
	VPRewards         []interface{}     `json:"vp_rewards,omitempty"`
	InventoryOptions  *InventoryOptions `json:"inventory_options,omitempty"`
	TotalContentPrice json.RawMessage   `json:"total_content_price,omitempty"`
	Content           []ContentItem     `json:"content,omitempty"`
}

type CatalogResponse struct {
	HasMore bool   `json:"has_more"`
	Items   []Item `json:"items"`
}
