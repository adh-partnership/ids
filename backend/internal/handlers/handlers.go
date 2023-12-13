package handlers

import (
	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/domain/charts"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	AirportHandler *AirportHandler
}

func RegisterHandlers(router chi.Router, airportService *airports.AirportService, chartService *charts.ChartService) (chi.Router, *Handlers) {
	h := &Handlers{}
	h.AirportHandler = NewAirportHandler(router, airportService)

	return router, h
}
