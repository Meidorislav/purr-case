package httpapi

import (
	"net/http"

	"github.com/go-chi/chi"
	middleware "github.com/go-chi/chi/v5/middleware"

	"purr-case/internal/httpapi/cases"
	"purr-case/internal/httpapi/global"
	"purr-case/internal/httpapi/inventory"
	"purr-case/internal/httpapi/items"
	"purr-case/internal/httpapi/payments"
	"purr-case/internal/httpapi/users"
)

func NewRouter(gh *global.Handler, uh *users.Handler, ih *items.Handler, ph *payments.Handler, invh *inventory.Handler, ch *cases.Handler) http.Handler {
	r := chi.NewRouter()

	// Middleware: Logging requests and recovering from panics.
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", gh.Health)

	// ---------------------------------------------------------------------------
	// Login
	// ---------------------------------------------------------------------------
	r.Group(func(r chi.Router) {
		r.Use(Auth)
		r.Get("/me", uh.Me)
	})

	// ---------------------------------------------------------------------------
	// Items
	// ---------------------------------------------------------------------------
	r.Route("/items", func(r chi.Router) {
		r.Use(OptionalAuth)
		r.Get("/", ih.GetItems)
		r.Get("/sku/{sku}", ih.GetItemBySku)
		r.Get("/virtual_items", ih.GetVirtualItems)
	})

	// ---------------------------------------------------------------------------
	// Inventory
	// ---------------------------------------------------------------------------
	r.Group(func(r chi.Router) {
		r.Use(Auth)
		r.Get("/inventory", invh.GetUserInventory)
		r.Post("/inventory/consume", invh.ConsumeInventoryItem) // POST /inventory/consume - consume a quantity of an item from the authenticated user's inventory.
		r.Get("/inventory/{sku}", invh.GetCurrencyQuantity)    // GET /inventory/{fish|food|yarn} - get quantity of a currency item.
	})

	// ---------------------------------------------------------------------------
	// Payments
	// ---------------------------------------------------------------------------
	r.Route("/payments", func(r chi.Router) { // group endpoints related to payments
		r.With(Auth).Post("/checkout", ph.CreateCheckout) // POST /payments/checkout - create a payment and get the payment link
		r.Post("/webhook", ph.HandleWebhook)              // POST /payments/webhook - handle webhook from the payment provider
	})

	// ---------------------------------------------------------------------------
	// Cases
	// ---------------------------------------------------------------------------
	r.Group(func(r chi.Router) {
		r.Use(Auth)
		r.Post("/cases/{sku}/open", ch.OpenCase)
	})

	return r
}
