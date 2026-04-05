package users

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	dto "purr-case/internal/dto/users"
	"purr-case/internal/httpapi/respond"
)

type Handler struct{}

func InitHandler() *Handler {
	return &Handler{}
}

const (
	userIdCtxKey = "userId"
	tokenCtxKey  = "token"
)

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(tokenCtxKey).(string)
	if !ok {
		respond.WriteError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://login.xsolla.com/api/users/me",
		nil,
	)
	if err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to create request")
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to get user info")
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		respond.WriteError(w, res.StatusCode, "failed to get user info")
		return
	}

	var userInfo dto.UserInfo
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to decode user info")
		return
	}

	respond.WriteJSON(w, http.StatusOK, userInfo)
}
