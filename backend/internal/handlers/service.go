package handlers

import (
	"net/http"

	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/ids/backend/pkg/response"
	"github.com/go-chi/chi/v5"
)

// This should be added by the main server package so it is not part of the versioned group
func NewServiceHandlers(router *chi.Mux) {
	logger.ZL.Info().Msg("Registering service handlers")
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		response.Respond(w, r, "pong", http.StatusOK)
	})
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		response.Respond(w, r, "ok", http.StatusOK)
	})
}
