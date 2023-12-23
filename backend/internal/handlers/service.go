package handlers

import (
	"net/http"

	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/ids/backend/pkg/response"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ServiceHandler struct {
	db *gorm.DB
}

// This should be added by the main server package so it is not part of the versioned group.
// These register routes that are used by Kubernetes or other health checkers to determine if the server is alive
// and should not be versioned.
func NewServiceHandler(router *chi.Mux, db *gorm.DB) *ServiceHandler {
	h := &ServiceHandler{
		db: db,
	}
	h.registerServiceHandlers(router)
	return h
}

func (h *ServiceHandler) registerServiceHandlers(router *chi.Mux) {
	logger.ZL.Info().Msg("Registering service handlers")
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		response.Respond(w, r, "pong", http.StatusOK)
	})
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		response.Respond(w, r, "ok", http.StatusOK)
	})
	router.Get("/readyz", h.Readyz)
}

func (h *ServiceHandler) Ping(w http.ResponseWriter, r *http.Request) {
	response.Respond(w, r, "pong", http.StatusOK)
}

func (h *ServiceHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	response.Respond(w, r, "ok", http.StatusOK)
}

func (h *ServiceHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	d, err := h.db.DB()
	if err != nil {
		response.Respond(w, r, "not ok", http.StatusInternalServerError)
		return
	}

	if err := d.Ping(); err != nil {
		response.Respond(w, r, "not ok", http.StatusInternalServerError)
		return
	}

	response.Respond(w, r, "ok", http.StatusOK)
}
