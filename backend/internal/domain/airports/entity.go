package airports

import "time"

type Airport struct {
	ID               int64      `json:"id"`
	FAAID            string     `json:"faa_id" gorm:"unique"`
	ICAOID           string     `json:"icao_id" gorm:"unique"`
	ATIS             string     `json:"atis"`
	ATISTime         *time.Time `json:"atis_time"`
	ArrivalATIS      string     `json:"arrival_atis"`
	ArrivalATISTime  *time.Time `json:"arrival_atis_time"`
	DepartureRunways string     `json:"departure_runways"`
	ArrivalRunways   string     `json:"arrival_runways"`
	METAR            string     `json:"metar"`
	TAF              string     `json:"taf"`
	MagVar           int        `json:"mag_var"`
	ParentFacility   int64      `json:"parent_facility"`
}
