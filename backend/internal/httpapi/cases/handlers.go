package cases

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	cases_dto "purr-case/internal/dto/cases"
	"purr-case/internal/httpapi/respond"
	cases_service "purr-case/internal/service/cases"
)

const userIdCtxKey = "userId"
const tokenCtxKey = "token"

type Handler struct {
	svc *cases_service.Service
}

func InitHandler(svc *cases_service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) OpenCase(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIdCtxKey).(string)
	if !ok || userID == "" {
		respond.WriteError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	caseSKU := chi.URLParam(r, "sku")
	if caseSKU == "" {
		respond.WriteError(w, http.StatusBadRequest, "missing sku")
		return
	}

	token, _ := r.Context().Value(tokenCtxKey).(string)

	wonSKU, err := h.svc.OpenCase(r.Context(), userID, caseSKU, token)
	if err != nil {
		switch {
		case errors.Is(err, cases_service.ErrNotInInventory):
			respond.WriteError(w, http.StatusConflict, "case not in inventory")
		case errors.Is(err, cases_service.ErrCaseNotFound):
			respond.WriteError(w, http.StatusNotFound, "case not found")
		case errors.Is(err, cases_service.ErrInvalidDropTable):
			respond.WriteError(w, http.StatusUnprocessableEntity, "case has no drop table")
		default:
			respond.WriteError(w, http.StatusInternalServerError, "failed to open case")
		}
		return
	}

	respond.WriteJSON(w, http.StatusOK, cases_dto.OpenCaseResponse{
		CaseSKU: caseSKU,
		WonSKU:  wonSKU,
	})
}
