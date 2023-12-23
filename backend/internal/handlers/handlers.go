package handlers

import (
	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/domain/auth"
	"github.com/adh-partnership/ids/backend/internal/domain/charts"
	"github.com/adh-partnership/ids/backend/internal/domain/pireps"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	AirportHandler *AirportHandler
	AuthHandler    *AuthHandler
	PIREPHandler   *PIREPHandler
}

type Services struct {
	AirportService *airports.AirportService
	AuthService    *auth.AuthService
	ChartService   *charts.ChartService
	PIREPService   *pireps.PIREPService
}

func RegisterHandlers(
	router chi.Router,
	services *Services,
) (chi.Router, *Handlers) {
	h := &Handlers{}
	h.AirportHandler = NewAirportHandler(router, services.AirportService)
	h.AuthHandler = NewAuthHandler(router, services.AuthService)
	h.PIREPHandler = NewPIREPHandler(router, services.PIREPService)

	return router, h
}
