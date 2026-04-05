package global

import (
	"net/http"
	"purr-case/internal/httpapi/respond"
)

type Handler struct{}

func InitHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	respond.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
