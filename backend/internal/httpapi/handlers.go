package httpapi

import (
	"encoding/json"
	"net/http"
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
