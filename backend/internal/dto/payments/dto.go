package payments_dto

// Structure for an item in the shopping cart
type CheckoutItem struct {
	SKU      string  `json:"sku"`      // Product identifier
	Name     string  `json:"name"`     // Product name
	Type     string  `json:"type"`     // Entity type: skin, case, battlepass
	Quantity int     `json:"quantity"` // Quantity of the product
	Price    float64 `json:"price"`    // Price per unit
	Currency string  `json:"currency"` // Currency (e.g., "USD", "EUR")
}

// Structure for a request to create a payment
type CreateCheckoutRequest struct {
	UserID string         `json:"userId"` // Temporary field 
	Items  []CheckoutItem `json:"items"`  // Contents of the cart
}

// Structure for the response when creating a payment
type CreateCheckoutResponse struct {
	OrderID     string `json:"orderId"`     // Order identifier
	Status      string `json:"status"`      // Payment status (e.g., "pending", "completed")
	ItemsCount  int    `json:"itemsCount"`  // Number of items in the cart
	CheckoutURL string `json:"checkoutUrl"` // Payment link
}
