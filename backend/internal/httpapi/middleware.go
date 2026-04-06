package httpapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"purr-case/internal/httpapi/respond"
	"strings"
	"time"
)

const (
	userIdCtxKey = "userId"
	tokenCtxKey  = "token"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respond.WriteError(w, http.StatusUnauthorized, "missing Authorization header")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			respond.WriteError(w, http.StatusUnauthorized, "invalid Authorization header")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		body := respond.MustJSON(map[string]string{"token": token})

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			"https://login.xsolla.com/api/token/validate",
			bytes.NewBuffer(body),
		)
		if err != nil {
			respond.WriteError(w, http.StatusInternalServerError, "failed to create request")
			return
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 204 {
			respond.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		userId, err := getIdFromJWT(token)
		if err != nil {
			respond.WriteError(w, http.StatusUnauthorized, "failed to parse token")
			return
		}

		ctx = context.WithValue(r.Context(), userIdCtxKey, userId)
		ctx = context.WithValue(ctx, tokenCtxKey, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getIdFromJWT(token string) (string, error) {
	payloadPart := strings.Split(token, ".")[1]
	payload, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return "", err
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", err
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	return userId, nil
}
