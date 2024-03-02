package handlers

import (
	"errors"
	"net/http"

	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/dtos"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/ids/backend/pkg/render"
	"github.com/adh-partnership/ids/backend/pkg/response"
	"github.com/adh-partnership/ids/backend/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type ExternalHandler struct {
	airportService *airports.AirportService
}

func NewExternalHandler(router chi.Router, airportService *airports.AirportService) *ExternalHandler {
	controller := &ExternalHandler{
		airportService: airportService,
	}

	router.Route("/external", func(r chi.Router) {
		r.Post("/vatis", controller.vatis)
	})

	return controller
}

// Soon(TM)
func (h *ExternalHandler) vatis(w http.ResponseWriter, r *http.Request) {
	var dto *dtos.VATISRequest
	if err := render.Bind(r, dto); err != nil {
		response.Respond(w, r, nil, http.StatusBadRequest)
		return
	}

	airport, err := h.airportService.GetAirport(dto.Facility)
	if err != nil || airport.FAAID == "" {
		if errors.Is(err, airports.ErrInvalidAirport) || airport.FAAID == "" {
			response.Respond(w, r, nil, http.StatusNotFound)
			return
		}
		logger.ZL.Error().Err(err).Msg("Error getting airport")
		response.Respond(w, r, nil, http.StatusInternalServerError)
		return
	}

	if dto.ATISType == "arrival" {
		airport.ArrivalATIS = dto.ATISLetter
		airport.ArrivalATISTime = utils.Now()
	} else {
		airport.ATIS = dto.ATISLetter
		airport.ATISTime = utils.Now()
	}

	if err := h.airportService.UpdateAirport(airport); err != nil {
		logger.ZL.Error().Err(err).Msg("Error updating airport")
		response.Respond(w, r, nil, http.StatusInternalServerError)
		return
	}

	response.Respond(w, r, airport.METAR, http.StatusOK)
}
