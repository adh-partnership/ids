package session

import (
	"net/http"

	gsess "github.com/gorilla/sessions"

	"github.com/adh-partnership/ids/backend/pkg/response"
)

func IsAuthenticated(r *http.Request) bool {
	if r.Context().Value("session") == nil {
		return false
	}

	session := r.Context().Value("session").(*gsess.Session)
	if session.Values["user"] != nil {
		return true
	}
	return false
}

func AuthenticatedMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			response.Respond(w, r, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
