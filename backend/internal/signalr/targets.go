package signalr

import (
	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/domain/pireps"
	"github.com/adh-partnership/ids/backend/internal/dtos"
)

// Invocation target names

const (
	RECEIVE_SERVER_MESSAGE = "ReceiveServerMessage"
	RECEIVE_AIRPORT_UPDATE = "ReceiveAirportUpdate"
	RECEIVE_PIREP          = "ReceivePirep"
)

// Helpers for sending messages to clients
func (h *IDSHub) SendServerMessage(connectionID string, message string) {
	h.Clients().Client(connectionID).Send(RECEIVE_SERVER_MESSAGE, message)
}

func (h *IDSHub) SendAirportUpdate(connectionID string, old, new airports.Airport) {
	h.Clients().Client(connectionID).Send(
		RECEIVE_AIRPORT_UPDATE,
		dtos.AirportResponseFromEntity(old),
		dtos.AirportResponseFromEntity(new),
	)
}

func (h *IDSHub) SendPirep(connectionID string, pirep *pireps.PIREP) {
	h.Clients().Client(connectionID).Send(
		RECEIVE_PIREP,
		pirep.ToString(),
	)
}
