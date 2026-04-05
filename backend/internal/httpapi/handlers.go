package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct{}

func InitHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(tokenCtxKey).(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid token")
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
		writeError(w, http.StatusInternalServerError, "failed to create request")
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get user info")
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		writeError(w, res.StatusCode, "failed to get user info")
		return
	}

	var userInfo UserInfo
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to decode user info")
		return
	}

	writeJSON(w, http.StatusOK, userInfo)
}

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

// writeJSON sets Content-Type: application/json and encodes body to JSON.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// writeError writes a JSON object {"error": "..."} with the given status code.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// mustJSON serializes a value to JSON. Panics on error (should not happen with DTOs).
func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic("json marshal: " + err.Error())
	}
	b = append(b, '\n') // for compatibility with json.Encoder, which adds \n
	return b
}
