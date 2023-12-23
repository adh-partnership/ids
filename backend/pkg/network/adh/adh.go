package adh

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/adh-partnership/ids/backend/pkg/config"
	"github.com/adh-partnership/ids/backend/pkg/network"
)

func GetUserInfo(cid string) (*User, error) {
	resp, err := network.Request("GET", fmt.Sprintf("%s%s", config.GetConfig().Facility.ADH.APIBase, fmt.Sprintf("/v1/user/%s", cid)), nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

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

func GetOverflights() ([]*Flightsv1, error) {
	resp, err := network.Request("GET", fmt.Sprintf(
		"%s/v1/overflight/%s",
		config.GetConfig().Facility.ADH.APIBase,
		config.GetConfig().Facility.Identifier,
	), nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var flights []*Flightsv1
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &flights)
	if err != nil {
		return nil, err
	}

	return flights, nil
}
