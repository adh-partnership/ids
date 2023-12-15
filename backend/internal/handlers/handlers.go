package handlers

import (
	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/domain/auth"
	"github.com/adh-partnership/ids/backend/internal/domain/charts"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	AirportHandler *AirportHandler
	AuthHandler    *AuthHandler
}

type Services struct {
	AirportService *airports.AirportService
	AuthService    *auth.AuthService
	ChartService   *charts.ChartService
}

func RegisterHandlers(
	router chi.Router,
	services *Services,
) (chi.Router, *Handlers) {
	h := &Handlers{}
	h.AirportHandler = NewAirportHandler(router, services.AirportService)
	h.AuthHandler = NewAuthHandler(router, services.AuthService)

	return router, h
}
