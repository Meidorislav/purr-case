package httpapi

import (
	"net/http"

	"github.com/go-chi/chi"
	middleware "github.com/go-chi/chi/v5/middleware"

	"purr-case/internal/httpapi/global"
	"purr-case/internal/httpapi/items"
	"purr-case/internal/httpapi/payments"
	"purr-case/internal/httpapi/users"
)

func NewRouter(gh *global.Handler, uh *users.Handler, ih *items.Handler, ph *payments.Handler) http.Handler {
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
		r.Get("/", ih.GetItems)
		r.Get("/{sku}", ih.GetItemBySku)
	})

	// ---------------------------------------------------------------------------
	// Payments
	// ---------------------------------------------------------------------------
	r.Route("/payments", func(r chi.Router) { // group endpoints related to payments
		r.Post("/checkout", ph.CreateCheckout) // POST /payments/checkout - create a payment and get the payment link
		r.Post("/webhook", ph.HandleWebhook)   // POST /payments/webhook - handle webhook from the payment provider
	})

	return r
}
