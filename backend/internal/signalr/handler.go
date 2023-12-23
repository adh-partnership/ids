package signalr

import (
	"net/http"
	"sync"

	"github.com/adh-partnership/ids/backend/internal/middleware/session"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/signalr"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type SignalRHandler struct {
	m       *signalr.HttpMux
	clients map[string]string
	mut     sync.RWMutex
}

var ()

func NewHandler(router chi.Router, server signalr.Server) *SignalRHandler {
	m := signalr.NewHTTPMux(server)
	h := &SignalRHandler{
		m:       m,
		clients: make(map[string]string),
	}

	router.Group(func(r chi.Router) {
		r.Use(session.AuthenticatedMiddleware)

		r.Get("/", h.HandleGet)
		r.Post("/negotiate", m.Negotiate)
	})

	logger.ZL.Debug().Msgf("h=%+v", h)

	return h
}

func (h *SignalRHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	sess := r.Context().Value("session").(*sessions.Session)
	if sess.Values["user"] == nil { // Should never happen since we have the middleware
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	connectionID := r.URL.Query().Get("id")
	if connectionID == "" {
		connectionID = uuid.NewString()
		q := r.URL.Query()
		q.Add("id", connectionID)
		r.URL.RawQuery = q.Encode()
		h.m.AddConnectionID(connectionID)
	}

	h.mut.Lock()
	h.clients[sess.Values["user"].(string)] = connectionID
	h.mut.Unlock()

	logger.ZL.Info().Msgf("SignalR connection %s for user %s", connectionID, sess.Values["user"].(string))
	logger.ZL.Debug().Msgf("SignalR clients: %+v", h.clients)

	h.m.ServeHTTP(w, r)
}

func (h *SignalRHandler) FindClientByConnectionId(connectionID string) string {
	h.mut.RLock()
	for k, v := range h.clients {
		if v == connectionID {
			h.mut.RUnlock()
			return k
		}
	}
	h.mut.RUnlock()
	return ""
}
