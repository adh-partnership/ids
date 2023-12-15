package session

import (
	"context"
	"net/http"
	"time"

	"github.com/adh-partnership/ids/backend/pkg/config"
	gsess "github.com/gorilla/sessions"
)

var (
	Store          *gsess.CookieStore
	session_config *config.Session
)

func New(sconfig *config.Session) {
	session_config = sconfig
	Store = gsess.NewCookieStore([]byte(session_config.HashKey), []byte(session_config.BlockKey))
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, session_config.Name)

		// Check if there is a user value in session, it should be a string if we set it
		session.Values["time"] = time.Now().Unix() // We do this to refresh the cookie
		session.Save(r, w)

		r = r.WithContext(context.WithValue(r.Context(), "session", session))

		h.ServeHTTP(w, r)
	})
}
