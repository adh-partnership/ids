package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/adh-partnership/ids/backend/pkg/config"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/ids/backend/pkg/network"
	"github.com/adh-partnership/ids/backend/pkg/network/adh"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	ErrInvalidProvider = errors.New("invalid provider")
)

type AuthService struct {
	OAuth2   *oauth2.Config
	Provider string
	ADHBase  string
	Rostered bool
	Store    sessions.Store
}

type AuthServiceSessions struct {
	r       *http.Request
	w       http.ResponseWriter
	Session *sessions.Session
}

func NewAuthService(store sessions.Store, oauth2provider string) *AuthService {
	c := config.GetConfig().OAuth
	oauth := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Scopes:       []string{"vatsim_details"},
		RedirectURL:  fmt.Sprintf("%s/v1/auth/callback", c.MyBaseURL),
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s%s", c.BaseURL, c.Endpoints.Authorization),
			TokenURL: fmt.Sprintf("%s%s", c.BaseURL, c.Endpoints.Token),
		},
	}

	oauth2provider = strings.ToLower(oauth2provider)
	if !config.IsValidOAuth2Provider(oauth2provider) {
		logger.ZL.Error().Err(ErrInvalidProvider).Msgf("invalid provider: %s", oauth2provider)
		return nil
	}

	return &AuthService{
		OAuth2:   oauth,
		Provider: oauth2provider,
		ADHBase:  config.GetConfig().Facility.ADH.APIBase,
		Rostered: config.GetConfig().Facility.ADH.Rostered,
		Store:    store,
	}
}

// Queries the userinfo endpoint with the provided oauth2 token. Will return
// the CID if available, whether or not the user is ok to use the system, and
// any errors that may have occurred.
//
// The criteria to determine if a user can use the system is:
// - Are they a VATSIM member?
// - If facility.adh.base_api is set and the facility.adh.rostered bool in config is true, are they rostered?
func (s *AuthService) GetUserInfo(token *oauth2.Token) (string, bool, error) {
	resp, err := network.Request("GET", fmt.Sprintf("%s%s", config.GetConfig().OAuth.BaseURL, config.GetConfig().OAuth.Endpoints.Userinfo), map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token.AccessToken),
	}, nil)
	if err != nil {
		return "", false, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, err
	}

	if resp.StatusCode != 200 {
		logger.ZL.Error().Msgf("received non-200 status code: %d", resp.StatusCode)
		logger.ZL.Debug().Msgf("Body: %s", string(contents))

		return "", false, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var cid string = ""
	if s.Provider == "vatsim" {
		data := &dtoVATSIMUserResponse{}
		err = json.Unmarshal(contents, data)
		if err != nil {
			logger.ZL.Error().Err(err).Msgf("failed to unmarshal vatsim user response: %s", string(contents))
			return "", false, err
		}
		cid = data.CID
	} else if s.Provider == "adh" {
		data := &dtoADHUserResponse{}
		err = json.Unmarshal(contents, data)
		if err != nil {
			logger.ZL.Error().Err(err).Msgf("failed to unmarshal adh user response: %s", string(contents))
			return "", false, err
		}

		if s.ADHBase != "" && s.Rostered {
			if data.User.ControllerType == "none" {
				return "", false, nil
			}
		}

		cid = fmt.Sprintf("%d", data.User.CID)
	} else {
		// We shouldn't make it here
		logger.ZL.Error().Err(ErrInvalidProvider).Msgf("invalid provider: %s", s.Provider)
		return "", false, ErrInvalidProvider
	}

	// We don't need to check ADH here since we can check it with the SSO response
	if s.Provider != "adh" && s.ADHBase != "" && s.Rostered {
		user, err := adh.GetUserInfo(cid)
		if err != nil {
			logger.ZL.Error().Err(err).Msgf("failed to get adh user info: %s", cid)
			return "", false, err
		}
		if user.ControllerType == "none" {
			return "", false, nil
		}
	}

	return cid, true, nil
}

func (s *AuthService) Session(r *http.Request, w http.ResponseWriter) *AuthServiceSessions {
	sess, _ := s.Store.Get(r, config.GetConfig().Session.Name)
	return &AuthServiceSessions{
		Session: sess,
		r:       r,
		w:       w,
	}
}

func (s *AuthServiceSessions) Save() {
	err := s.Session.Save(s.r, s.w)
	if err != nil {
		logger.ZL.Error().Err(err).Msgf("failed to save session: %+v", err)
	}
}

func (s *AuthServiceSessions) Get(key string) interface{} {
	return s.Session.Values[key]
}

func (s *AuthServiceSessions) Set(key string, value interface{}) {
	s.Session.Values[key] = value
}

func (s *AuthServiceSessions) Delete(key string) {
	delete(s.Session.Values, key)
}
