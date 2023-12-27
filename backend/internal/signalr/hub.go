package signalr

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/domain/pireps"
	"github.com/adh-partnership/ids/backend/internal/dtos"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/signalr"
	"github.com/go-chi/chi/v5"
)

type IDSHub struct {
	signalr.Hub
	server signalr.Server
}

var (
	handler        *SignalRHandler
	hub            IDSHub
	airportService *airports.AirportService
	pirepService   *pireps.PIREPService
)

func New(ctx context.Context, r chi.Router, as *airports.AirportService, ps *pireps.PIREPService) (*IDSHub, error) {
	hub = IDSHub{}
	l := &Logger{}

	server, err := signalr.NewServer(ctx,
		signalr.SimpleHubFactory(&hub),
		signalr.HTTPTransports(signalr.TransportWebSockets), // Only allow websocket connections
		signalr.AllowOriginPatterns([]string{"*"}),
		signalr.KeepAliveInterval(5*time.Second),
		signalr.Logger(l, false),
	)
	if err != nil {
		return nil, err
	}

	hub.server = server
	airportService = as
	pirepService = ps

	handler = NewHandler(r, server)

	return &hub, nil
}

func (h *IDSHub) ConfigureHooks(airportService *airports.AirportService, pirepService *pireps.PIREPService) {
	airportService.AddHook(AirportHook)
	pirepService.AddHook(PIREPHook)
}

// We setup a clients table in our Handler which tracks connectionID to CID already
func (h *IDSHub) OnConnected(connectionID string) {
	h.Groups().AddToGroup("sessions", connectionID)
	airports, err := airportService.GetAirports()
	if err != nil {
		logger.ZL.Error().Err(err).Msg("Error getting airports")
		return
	}
	h.Clients().Client(connectionID).Send("airports", dtos.AirportResponsesFromEntities(airports))
}

func (h *IDSHub) OnDisconnected(connectionID string) {
	if handler == nil {
		logger.ZL.Error().Msg("SignalR handler not initialized ?!")
		return
	}
	user := handler.FindClientByConnectionId(connectionID)
	logger.ZL.Info().Msgf(
		"SignalR connection %s disconnected, user: %s",
		connectionID,
		user,
	)

	if user != "" {
		handler.mut.Lock()
		delete(handler.clients, user)
		handler.mut.Unlock()
	}
	h.Groups().RemoveFromGroup("sessions", connectionID)
}

func (h *IDSHub) UpdateAirport(id string, patch json.RawMessage) json.RawMessage {
	var airportPatch dtos.AirportPatch
	if err := json.Unmarshal(patch, &airportPatch); err != nil {
		panic(err)
	}

	airport, err := airportService.GetAirport(id)
	if err != nil {
		panic(err)
	}

	airportPatch.MergeInto(&airport)
	if err := airportService.UpdateAirport(airport); err != nil {
		panic(err)
	}

	data, err := json.Marshal(dtos.AirportResponseFromEntity(airport))
	if err != nil {
		panic(err)
	}

	return data
}

func (h *IDSHub) SubmitPIREP(pirep json.RawMessage) {
	var pirepDTO dtos.PIREPRequest
	if err := json.Unmarshal(pirep, &pirepDTO); err != nil {
		panic(err)
	}

	err := pirepService.CreatePIREP(
		pirepDTO.ToEntity(),
	)
	if err != nil {
		panic(err)
	}
}

func AirportHook(old, new airports.Airport) {
	hub.server.HubClients().All().Send(
		"airportUpdate",
		dtos.AirportResponseFromEntity(old),
		dtos.AirportResponseFromEntity(new),
	)
}

func PIREPHook(pirep *pireps.PIREP) {
	hub.server.HubClients().All().Send(
		"pirepUpdate",
		dtos.PIREPResponseFromEntity(pirep),
	)
}
