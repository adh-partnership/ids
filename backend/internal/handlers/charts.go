package handlers

import (
	"errors"
	"net/http"

	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/domain/charts"
	"github.com/adh-partnership/ids/backend/internal/dtos"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/ids/backend/pkg/response"
	"github.com/go-chi/chi/v5"
)

type ChartHandler struct {
	AirportService *airports.AirportService
	ChartService   *charts.ChartService
}

func NewChartHandler(router chi.Router, chartService *charts.ChartService, airportService *airports.AirportService) *ChartHandler {
	controller := &ChartHandler{
		ChartService:   chartService,
		AirportService: airportService,
	}

	router.Route("/charts", func(r chi.Router) {
		r.Get("/", controller.GetAllCharts)
		r.Get("/{id}", controller.GetCharts)
	})

	return controller
}

func (h *ChartHandler) GetAllCharts(w http.ResponseWriter, r *http.Request) {
	var c map[string][]*charts.Chart
	var err error

	c, err = h.ChartService.GetAllCharts()
	if err != nil {
		logger.ZL.Error().Err(err).Msg("Error getting charts")
		response.Respond(w, r, nil, http.StatusInternalServerError)
		return
	}

	response.Respond(w, r, dtos.ChartSeparatedResponsesFromEntities(h.AirportService, c), http.StatusOK)
}

func (h *ChartHandler) GetCharts(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var c []*charts.Chart
	var err error

	c, err = h.ChartService.GetCharts(id)

	if err != nil {
		if errors.Is(err, charts.ErrInvalidAiport) || errors.Is(err, charts.ErrNoCharts) {
			response.Respond(w, r, nil, http.StatusNotFound)
			return
		}
		logger.ZL.Error().Err(err).Msg("Error getting charts")
		response.Respond(w, r, nil, http.StatusInternalServerError)
		return
	}

	response.Respond(w, r, dtos.ChartResponsesFromEntities(c), http.StatusOK)
}
