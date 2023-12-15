// Typically DTOs would go to the DTOs package, however, these are only used internally by this service, so it's fine
// to leave them here

package auth

// We only care about the CID
type dtoVATSIMUserResponse struct {
	CID string `json:"cid"`
}

// When the API gets its rewrite, we can import its DTO instead of writing our own
type dtoADHUserResponse struct {
	User struct {
		CID            uint   `json:"cid"`
		ControllerType string `json:"controller_type"`
	} `json:"user"`
}
