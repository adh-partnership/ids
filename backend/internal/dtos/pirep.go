package dtos

import (
	"errors"
	"net/http"
	"time"

	"github.com/adh-partnership/ids/backend/internal/domain/pireps"
)

type PIREPRequest struct {
	Type string     `json:"type"`
	OV   string     `json:"OV"`
	TM   *time.Time `json:"TM"`
	FL   string     `json:"FL"`
	TP   string     `json:"TP"`
	SK   string     `json:"SK"`
	WX   string     `json:"WX"`
	TA   string     `json:"TA"`
	WV   string     `json:"WV"`
	TB   string     `json:"TB"`
	IC   string     `json:"IC"`
	RM   string     `json:"RM"`
}

func (p *PIREPRequest) Bind(r *http.Request) error {
	if p.Type != "UA" && p.Type != "UUA" {
		return errors.New("invalid type")
	}

	if p.OV == "" {
		return errors.New("missing OV field")
	}

	if p.TM == nil {
		return errors.New("missing TM field")
	}

	return nil
}

func (p *PIREPRequest) ToEntity() *pireps.PIREP {
	pirep := &pireps.PIREP{
		Type:         p.Type,
		Location:     p.OV,
		TM:           p.TM,
		Altitude:     p.FL,
		AircraftType: p.TP,
		SkyCover:     p.SK,
		Visibility:   p.WX,
		Temperature:  p.TA,
		Wind:         p.WV,
		Turbulence:   p.TB,
		Icing:        p.IC,
		Remarks:      p.RM,
	}
	pirep.Raw = pirep.ToString()
	return pirep
}

type PIREPResponse struct {
	ID   int64      `json:"id,omitempty"`
	Type string     `json:"type,omitempty"`
	OV   string     `json:"OV,omitempty"`
	TM   *time.Time `json:"TM,omitempty"`
	FL   string     `json:"FL,omitempty"`
	TP   string     `json:"TP,omitempty"`
	SK   string     `json:"SK,omitempty"`
	WX   string     `json:"WX,omitempty"`
	TA   string     `json:"TA,omitempty"`
	WV   string     `json:"WV,omitempty"`
	TB   string     `json:"TB,omitempty"`
	IC   string     `json:"IC,omitempty"`
	RM   string     `json:"RM,omitempty"`
	Raw  string     `json:"raw,omitempty"`
}

func NewPIREPResponse(pirep *pireps.PIREP) *PIREPResponse {
	p := &PIREPResponse{
		ID: pirep.ID,
		OV: pirep.Location,
		TM: pirep.TM,
		FL: pirep.Altitude,
		TP: pirep.AircraftType,
		SK: pirep.SkyCover,
		WX: pirep.Visibility,
		TA: pirep.Temperature,
		WV: pirep.Wind,
		TB: pirep.Turbulence,
		IC: pirep.Icing,
		RM: pirep.Remarks,
	}

	if p.Raw == "" {
		p.Raw = pirep.ToString()
	}
	return p
}

func NewPIREPResponses(pireps []*pireps.PIREP) []*PIREPResponse {
	var pirepResponses []*PIREPResponse
	for _, pirep := range pireps {
		pirepResponses = append(pirepResponses, NewPIREPResponse(pirep))
	}
	return pirepResponses
}

func PIREPResponseFromEntity(pirep *pireps.PIREP) *PIREPResponse {
	return NewPIREPResponse(pirep)
}
