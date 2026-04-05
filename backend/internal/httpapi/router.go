package httpapi

import (
	"net/http"

	"github.com/go-chi/chi"
	middleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	// Middleware: Logging requests and recovering from panics.
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", h.Health)

	r.Route("/payments", func(r chi.Router){ // группируем эндпоинты, связанные с платежами
		r.Post("/checkout", h.CreateCheckout) // POST /payments/checkout - создать платеж и получить ссылку на оплату
		r.Post("/webhook", h.HandleWebhook)   // POST /payments/webhook - обработать вебхук от платежного провайдера
	})

	return r
}
