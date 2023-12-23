package adh

import "time"

// @TODO Replace this when the API is rewritten with an imported DTO
type User struct {
	CID                  uint                           `json:"cid" yaml:"cid" xml:"cid"`
	FirstName            string                         `json:"first_name" yaml:"first_name" xml:"first_name"`
	LastName             string                         `json:"last_name" yaml:"last_name" xml:"last_name"`
	OperatingInitials    string                         `json:"operating_initials" yaml:"operating_initials" xml:"operating_initials"`
	ControllerType       string                         `json:"controller_type" yaml:"controller_type" xml:"controller_type"`
	Certifications       map[string]*UserCertifications `json:"certifications" yaml:"certifications" xml:"certifications"`
	Rating               string                         `json:"rating" yaml:"rating" xml:"rating"`
	Status               string                         `json:"status" yaml:"status" xml:"status"`
	Roles                []string                       `json:"roles" yaml:"roles" xml:"roles"`
	Region               string                         `json:"region" yaml:"region" xml:"region"`
	Division             string                         `json:"division" yaml:"division" xml:"division"`
	Subdivision          string                         `json:"subdivision" yaml:"subdivision" xml:"subdivision"`
	DiscordID            string                         `json:"discord_id" yaml:"discord_id" xml:"discord_id"`
	RosterJoinDate       string                         `json:"roster_join_date" yaml:"roster_join_date" xml:"roster_join_date"`
	ExemptedFromActivity *bool                          `json:"exempted_from_activity" yaml:"exempted_from_activity" xml:"exempted_from_activity"`
	CreatedAt            string                         `json:"created_at" yaml:"created_at" xml:"created_at"`
	UpdatedAt            string                         `json:"updated_at" yaml:"updated_at" xml:"updated_at"`
}

type UserCertifications struct {
	DisplayName string `json:"display_name" yaml:"display_name" xml:"display_name"`
	Value       string `json:"value" yaml:"value" xml:"value"`
	Order       uint   `json:"order" yaml:"order" xml:"order"`
	Hidden      bool   `json:"hidden" yaml:"hidden" xml:"hidden"`
}

type Flightsv1 struct {
	Callsign    string     `json:"callsign" example:"N462AW"`
	CID         int        `json:"cid" example:"876594"`
	Facility    string     `json:"facility" example:"ZDV"`
	Latitude    float32    `json:"lat" example:"-33.867"`
	Longitude   float32    `json:"lon" example:"151.206"`
	Groundspeed int        `json:"spd" example:"150"`
	Heading     int        `json:"hdg" example:"180"`
	Altitude    int        `json:"alt" example:"10500"`
	Aircraft    string     `json:"type" example:"C208"`
	Departure   string     `json:"dep" example:"KLMO"`
	Arrival     string     `json:"arr" example:"KLMO"`
	Route       string     `json:"route" example:"DCT"`
	UpdatedAt   *time.Time `json:"lastSeen" example:"2020-01-01T00:00:00Z"`
}
