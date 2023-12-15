package dtos

import "net/http"

type VATISRequest struct {
	Facility          string `json:"facility"`
	Preset            string `json:"preset"`
	ATISLetter        string `json:"atisLetter"`
	ATISType          string `json:"atisType"`
	AirportConditions string `json:"airportConditions"`
	NOTAMs            string `json:"notams"`
	Timestamp         string `json:"timestamp"`
	Version           string `json:"version"`
}

func (a *VATISRequest) Bind(r *http.Request) error {
	return nil
}
