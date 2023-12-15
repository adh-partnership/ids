package signalr

import "net/http"

type NegotiateResponse struct {
	ConnectionID        string                       `json:"connectionId:omitempty"`
	ConnectionToken     string                       `json:"connectionToken:omitempty"`
	NegotiateVersion    int                          `json:"negotiateVersion:omitempty"`
	AvailableTransports []NegotiateResponseTransport `json:"availableTransports:omitempty"`
	Error               string                       `json:"error:omitempty"`
}

type NegotiateResponseTransport struct {
	Transport       string   `json:"transport:omitempty"`
	TransferFormats []string `json:"transferFormats:omitempty"`
}

type NegotiationRequest struct {
	NegotiateVersion int `json:"negotiateVersion"`
}

func (n *NegotiationRequest) Bind(r *http.Request) error {
	return nil
}
