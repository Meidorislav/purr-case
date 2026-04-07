package cases_dto

import "encoding/json"

type DropEntry struct {
	SKU    string `json:"sku"`
	Weight int    `json:"weight"`
}

type CaseCustomAttributes struct {
	Type      string      `json:"type"`
	DropTable []DropEntry `json:"drop_table"`
}

type XsollaItem struct {
	ItemID          int             `json:"item_id"`
	SKU             string          `json:"sku"`
	Type            string          `json:"type"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	ImageURL        *string         `json:"image_url"`
	Price           json.RawMessage `json:"price"`
	VirtualPrices   json.RawMessage `json:"virtual_prices"`
	Attributes      json.RawMessage `json:"attributes"`
	CustomAttributes json.RawMessage `json:"custom_attributes"`
	IsFree          bool            `json:"is_free"`
	Groups          json.RawMessage `json:"groups"`
	VirtualItemType string          `json:"virtual_item_type"`
	InventoryOptions json.RawMessage `json:"inventory_options"`
	Quantity        int             `json:"quantity"`
	Limits          json.RawMessage `json:"limits"`
}

type OpenCaseResponse struct {
	CaseSKU string     `json:"case_sku"`
	WonItem XsollaItem `json:"won_item"`
}
