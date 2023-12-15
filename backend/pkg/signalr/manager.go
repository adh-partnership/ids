package signalr

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/adh-partnership/ids/backend/pkg/render"
	"github.com/adh-partnership/ids/backend/pkg/response"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	webSocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (m *Manager) Negotiate(w http.ResponseWriter, r *http.Request) {
	// Check if negotiateVersion is in body
	dto := &NegotiationRequest{}
	var negVersion int = -1
	if err := render.Bind(r, dto); err == nil {
		negVersion = dto.NegotiateVersion
	} else {
		// Check if negotiateVersion is in query
		negVersion, err = strconv.Atoi(r.URL.Query().Get("negotiateVersion"))
		if err != nil {
			negVersion = 0
		}
	}

	if negVersion > 1 {
		response.Respond(w, r, &NegotiateResponse{
			Error: "This connection is not allowed.",
		}, http.StatusBadRequest)
	}

	var connectionToken string
	if negVersion == 1 {
		connectionToken = GenerateRandomString(32)
	}

	response.Respond(w, r, &NegotiateResponse{
		ConnectionID:     uuid.NewString(),
		ConnectionToken:  connectionToken,
		NegotiateVersion: negVersion,
		AvailableTransports: []NegotiateResponseTransport{
			{
				Transport: "WebSockets",
				TransferFormats: []string{
					"Text", // We could *probably* support Binary, too... but nah.
				},
			},
		},
	}, http.StatusOK)
}
