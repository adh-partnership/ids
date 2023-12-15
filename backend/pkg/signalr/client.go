package signalr

type ClientState int

const (
	ClientCreated ClientState = iota
	ClientConnecting
	ClientConnected
	ClientClosed
)

type Client struct {
	ConnectionID    string      `json:"connectionId"`
	ConnectionToken string      `json:"connectionToken"`
	State           ClientState `json:"-"` // Internal use only
}

type ClientList map[*Client]bool
