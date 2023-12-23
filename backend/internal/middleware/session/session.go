package session

import (
	"context"
	"net/http"
	"time"

	"github.com/adh-partnership/ids/backend/pkg/config"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	gsess "github.com/gorilla/sessions"
)

var (
	Store          *gsess.CookieStore
	session_config *config.Session
)

func New(sconfig *config.Session) {
	session_config = sconfig
	logger.ZL.Info().Msg("Configuring session store")
	Store = gsess.NewCookieStore([]byte(session_config.HashKey), []byte(session_config.BlockKey))
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.ZL.Info().Msg("Handling session middleware in request")
		session, err := Store.Get(r, session_config.Name)
		if err != nil {
			logger.ZL.Error().Err(err).Msg("unable to get session")
		}

		// Check if there is a user value in session, it should be a string if we set it
		session.Values["time"] = time.Now().Unix() // We do this to refresh the cookie
		err = session.Save(r, w)
		if err != nil {
			logger.ZL.Error().Err(err).Msg("unable to save session")
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))

		h.ServeHTTP(w, r)
	})
}
