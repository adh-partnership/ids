package handlers

import (
	"net/http"

	"github.com/adh-partnership/ids/backend/internal/domain/auth"
	"github.com/adh-partnership/ids/backend/internal/middleware/session"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/ids/backend/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(router chi.Router, authService *auth.AuthService) *AuthHandler {
	controller := &AuthHandler{
		authService: authService,
	}

	router.Route("/auth", func(r chi.Router) {
		r.Get("/login", controller.Login)
		r.Get("/callback", controller.Callback)
		r.Get("/logout", controller.Logout)
		r.Get("/check", controller.Check)
	})

	return controller
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state := uuid.NewString()
	verifier := oauth2.GenerateVerifier()

	url := h.authService.OAuth2.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(verifier),
	)

	s := h.authService.Session(r, w)
	s.Set("state", state)
	s.Set("verifier", verifier)
	s.Set("redirect", r.URL.Query().Get("redirect"))
	s.Save()

	response.Redirect(w, r, url, http.StatusFound)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	s := h.authService.Session(r, w)
	s.Delete("state")
	s.Delete("verifier")
	s.Delete("user")
	s.Save()

	response.Respond(w, r, nil, http.StatusOK)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	s := h.authService.Session(r, w)

	if s == nil {
		logger.ZL.Error().Msgf("failed to get session: %+v, %+v", s, session.Store)
		response.Respond(w, r, "Error", http.StatusInternalServerError)
		return
	}

	verifier := s.Get("verifier").(string)
	state := s.Get("state").(string)
	redirect := s.Get("redirect").(string)
	s.Delete("state")
	s.Delete("verifier")
	s.Delete("redirect")
	s.Save()

	if r.URL.Query().Get("state") != state && state != "" {
		logger.ZL.Debug().Msg("invalid oauth state from client")
		response.Respond(w, r, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := h.authService.OAuth2.Exchange(r.Context(), r.URL.Query().Get("code"), oauth2.VerifierOption(verifier))
	if err != nil {
		logger.ZL.Error().Err(err).Msg("failed to exchange token")
		response.Respond(w, r, "Error", http.StatusInternalServerError)
		return
	}

	cid, ok, err := h.authService.GetUserInfo(token)
	if err != nil {
		logger.ZL.Error().Err(err).Msg("failed to get user info")
		response.Respond(w, r, "Error", http.StatusInternalServerError)
		return
	}
	if !ok {
		logger.ZL.Debug().Msg("user not rostered")
		response.Respond(w, r, "Forbidden", http.StatusUnauthorized)
		return
	}

	s.Set("user", cid)
	s.Save()

	if redirect != "" {
		response.Redirect(w, r, redirect, http.StatusFound)
		return
	}

	response.Respond(w, r, "Success", http.StatusOK)
}

func (h *AuthHandler) Check(w http.ResponseWriter, r *http.Request) {
	s := h.authService.Session(r, w)
	if s.Get("user") == nil {
		response.Respond(w, r, nil, http.StatusUnauthorized)
		return
	}

	response.Respond(w, r, nil, http.StatusOK)
}
