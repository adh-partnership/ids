package adh

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/adh-partnership/ids/backend/pkg/config"
	"github.com/adh-partnership/ids/backend/pkg/network"
)

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

func GetUserInfo(cid string) (*User, error) {
	resp, err := network.Request("GET", fmt.Sprintf("%s%s", config.GetConfig().Facility.ADH.APIBase, fmt.Sprintf("/v1/users/%s", cid)), nil, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var user *User
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
