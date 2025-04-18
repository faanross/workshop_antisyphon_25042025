package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const AgentUUIDKey contextKey = "agentUUID"

func UUIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		agentUUID := r.Header.Get("X-Agent-ID")

		ctx := context.WithValue(r.Context(), AgentUUIDKey, agentUUID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
