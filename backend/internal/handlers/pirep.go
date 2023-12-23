package handlers

import (
	"net/http"

	"github.com/adh-partnership/ids/backend/internal/domain/pireps"
	"github.com/adh-partnership/ids/backend/internal/dtos"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/ids/backend/pkg/render"
	"github.com/adh-partnership/ids/backend/pkg/response"
	"github.com/go-chi/chi/v5"
)

type PIREPHandler struct {
	PIREPService *pireps.PIREPService
}

func NewPIREPHandler(router chi.Router, pirepService *pireps.PIREPService) *PIREPHandler {
	controller := &PIREPHandler{PIREPService: pirepService}

	router.Route("/pireps", func(r chi.Router) {
		r.Get("/", controller.GetPIREPs)
		r.Put("/", controller.PutPIREP)
	})

	return controller
}

func (h *PIREPHandler) GetPIREPs(w http.ResponseWriter, r *http.Request) {
	pireps, err := h.PIREPService.GetPIREPs()
	if err != nil {
		logger.ZL.Error().Err(err).Msg("Error getting pireps")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Respond(w, r, dtos.NewPIREPResponses(pireps), http.StatusOK)
}

func (h *PIREPHandler) PutPIREP(w http.ResponseWriter, r *http.Request) {
	pirep := &dtos.PIREPRequest{}
	if err := render.Bind(r, pirep); err != nil {
		logger.ZL.Error().Err(err).Msg("Error binding pirep")
		response.Respond(w, r, err, http.StatusBadRequest)
		return
	}

	if err := h.PIREPService.CreatePIREP(
		pirep.ToEntity(),
	); err != nil {
		logger.ZL.Error().Err(err).Msg("Error adding pirep")
		response.Respond(w, r, nil, http.StatusInternalServerError)
		return
	}

	response.Respond(
		w, r,
		nil,
		http.StatusNoContent,
	)
}
