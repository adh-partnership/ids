package handlers

import (
	"errors"
	"net/http"

	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/dtos"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/ids/backend/pkg/render"
	"github.com/adh-partnership/ids/backend/pkg/response"
	"github.com/go-chi/chi/v5"
)

type AirportHandler struct {
	AirportService *airports.AirportService
}

func NewAirportHandler(router chi.Router, airportService *airports.AirportService) *AirportHandler {
	controller := &AirportHandler{AirportService: airportService}

	router.Route("/airports", func(r chi.Router) {
		r.Get("/", controller.GetAirports)
		r.Get("/{id}", controller.GetAirport)
		r.Patch("/{id}", controller.PatchAirport)
	})

	return controller
}

func (h *AirportHandler) GetAirports(w http.ResponseWriter, r *http.Request) {
	airports, err := h.AirportService.GetAirports()
	if err != nil {
		logger.ZL.Error().Err(err).Msg("Error getting airports")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Respond(w, r, airports, http.StatusOK)
}

func (h *AirportHandler) GetAirport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	airport, err := h.AirportService.GetAirport(id)
	if err != nil {
		if errors.Is(err, airports.ErrInvalidAirport) {
			response.Respond(w, r, nil, http.StatusNotFound)
			return
		}
		logger.ZL.Error().Err(err).Msg("Error getting airport")
		response.Respond(w, r, nil, http.StatusInternalServerError)
		return
	}

	response.Respond(w, r, airport, http.StatusOK)
}

func (h *AirportHandler) PatchAirport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	airport, err := h.AirportService.GetAirport(id)
	if err != nil {
		if errors.Is(err, airports.ErrInvalidAirport) {
			response.Respond(w, r, nil, http.StatusNotFound)
			return
		}
		logger.ZL.Error().Err(err).Msg("Error getting airport")
		response.Respond(w, r, nil, http.StatusInternalServerError)
		return
	}

	patch := &dtos.AirportPatch{}
	if err := render.Bind(r, patch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	patch.MergeInto(airport)
	h.AirportService.UpdateAirport(airport)

	response.Respond(w, r, airport, http.StatusOK)
}
