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

	logger.ZL.Info().Msg("Registering airport handlers")
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

	logger.ZL.Info().Msgf("Getting airport %s", id)
	airport, err := h.AirportService.GetAirport(id)
	if err != nil {
		if errors.Is(err, airports.ErrInvalidAirport) {
			response.Respond(w, r, nil, http.StatusNotFound)
			return
		}
		response.Respond(w, r, nil, http.StatusInternalServerError)
		return
	}

	logger.ZL.Info().Msgf("Parsing patch for airport %s", airport.FAAID)
	patch := &dtos.AirportPatch{}
	if err := render.Bind(r, patch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.ZL.Info().Msgf("Patching airport %s", airport.FAAID)
	patch.MergeInto(airport)
	logger.ZL.Info().Msgf("Updating airport %s", airport.FAAID)
	h.AirportService.UpdateAirport(airport)
	logger.ZL.Info().Msgf("Done, responding to client")

	response.Respond(w, r, airport, http.StatusOK)
}
