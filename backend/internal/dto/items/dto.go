package items_dto

type Price struct {
	Amount                string `json:"amount"`
	AmountWithoutDiscount string `json:"amount_without_discount"`
	Currency              string `json:"currency"`
}

type VirtualPrice struct {
	ItemID                int     `json:"item_id"`
	SKU                   string  `json:"sku"`
	Name                  string  `json:"name"`
	Type                  string  `json:"type"`
	Description           string  `json:"description"`
	ImageURL              *string `json:"image_url"`
	Amount                int     `json:"amount"`
	AmountWithoutDiscount int     `json:"amount_without_discount"`
	IsDefault             bool    `json:"is_default"`
}

type Group struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}

type Period struct {
	DateFrom  string `json:"date_from"`
	DateUntil string `json:"date_until"`
}

type Item struct {
	ItemID          int            `json:"item_id"`
	SKU             string         `json:"sku"`
	Type            string         `json:"type"`
	VirtualItemType string         `json:"virtual_item_type"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	ImageURL        *string        `json:"image_url"`
	Price           *Price         `json:"price"`
	VirtualPrices   []VirtualPrice `json:"virtual_prices"`
	CanBeBought     bool           `json:"can_be_bought"`
	IsFree          bool           `json:"is_free"`
	Groups          []Group        `json:"groups"`
	Periods         []Period       `json:"periods"`
}

type CatalogResponse struct {
	HasMore bool   `json:"has_more"`
	Items   []Item `json:"items"`
}
