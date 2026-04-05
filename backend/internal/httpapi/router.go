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

	return r
}
