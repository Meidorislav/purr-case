package payments_dto

// CheckoutItem describes a single cart line that will later be sent to Xsolla.
// We only trust SKU and quantity from the client. Price, currency and all
// promotions must come from Xsolla Catalog / order APIs.
type CheckoutItem struct {
	SKU      string `json:"sku"`      // Product SKU from Xsolla Catalog.
	Quantity int    `json:"quantity"` // Requested quantity for this SKU.
}

// CreateCheckoutRequest describes the cart payload accepted by the checkout API.
type CreateCheckoutRequest struct {
	Items []CheckoutItem `json:"items"` // Cart contents to be converted into an order.
}

// CreateCheckoutResponse describes the current checkout session summary.
// The backend creates an Xsolla order/payment token server-side and returns
// a Pay Station URL built from the received token.
type CreateCheckoutResponse struct {
	OrderID     string `json:"orderId"`     // Local order identifier.
	Status      string `json:"status"`      // Order status. "new" is the closest Xsolla order state.
	ItemsCount  int    `json:"itemsCount"`  // Number of cart lines in the checkout session.
	CheckoutURL string `json:"checkoutUrl"` // Pay Station URL built from a payment token.
}
