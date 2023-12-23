package pireps

import (
	"time"
)

type PIREP struct {
	ID           int64      `json:"id"`
	Type         string     `json:"type"`
	Location     string     `json:"location"`
	TM           *time.Time `json:"time"` // time is a reserved word.. so use this
	Altitude     string     `json:"altitude"`
	AircraftType string     `json:"aircraft_type"`
	SkyCover     string     `json:"sky_cover"`
	Visibility   string     `json:"visibility"`
	Temperature  string     `json:"temperature"`
	Wind         string     `json:"wind"`
	Turbulence   string     `json:"turbulence"`
	Icing        string     `json:"icing"`
	Remarks      string     `json:"remarks"`
	Raw          string     `json:"raw"`
}

func (p *PIREP) ToString() string {
	str := p.Type + " "

	if p.Location != "" {
		str += "/OV " + p.Location
	}

	if p.TM != nil {
		str += "/TM " + p.TM.Format("1504")
	}

	if p.Altitude != "" {
		str += "/FL " + p.Altitude
	}

	if p.AircraftType != "" {
		str += "/TP " + p.AircraftType
	}

	if p.SkyCover != "" {
		str += "/SK " + p.SkyCover
	}

	if p.Visibility != "" {
		str += "/WX " + p.Visibility
	}

	if p.Temperature != "" {
		str += "/TA " + p.Temperature
	}

	if p.Wind != "" {
		str += "/WV " + p.Wind
	}

	if p.Turbulence != "" {
		str += "/TB " + p.Turbulence
	}

	if p.Icing != "" {
		str += "/IC " + p.Icing
	}

	if p.Remarks != "" {
		str += "/RM " + p.Remarks
	}

	return str
}
