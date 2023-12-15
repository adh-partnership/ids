package config

import (
	"errors"
	"strings"
)

var (
	ErrInvalidCacheDriver    = errors.New("invalid cache driver")
	ErrInvalidCacheRedisHost = errors.New("invalid redis host or port")
	ErrInvalidCacheRedisDB   = errors.New("invalid redis db")

	ErrInvalidOAuthProvider = errors.New("invalid oauth provider")

	ErrInvalidSessionBlockKey = errors.New("invalid session block key")
	ErrInvalidSessionHashKey  = errors.New("invalid session hash key")
	ErrMissingSessionDomain   = errors.New("missing session domain")
)

func ValidateConfig(c *Config) error {
	if err := ValidateCache(&c.Cache); err != nil {
		return err
	}

	if err := ValidateSession(&c.Session); err != nil {
		return err
	}

	return nil
}

func ValidateCache(c *Cache) error {
	if c.Driver == "redis" {
		if c.Host == "" || c.Port == 0 {
			return ErrInvalidCacheRedisHost
		}
		if c.DB < 0 || c.DB > 15 {
			return ErrInvalidCacheRedisDB
		}
	}

	if c.Driver != "redis" && c.Driver != "memory" {
		return ErrInvalidCacheDriver
	}

	if c.DefaultExpiration == nil {
		c.DefaultExpiration = &CacheExpiration{}
	}
	if c.DefaultExpiration.Airports == 0 {
		c.DefaultExpiration.Airports = 5 * 60
	}

	if c.DefaultExpiration.Charts == 0 {
		c.DefaultExpiration.Charts = 60 * 60
	}

	return nil
}

func ValidateOAuth(o *OAuth) error {
	if !IsValidOAuth2Provider(o.Provider) {
		return ErrInvalidOAuthProvider
	}

	return nil
}

func ValidateSession(s *Session) error {
	// Technically can be 16, 24, or 32 bytes... but force AES-256 with 32 bytes.
	// Future proof, allow more than 32 bytes but still in intervals of 8.
	if len(s.BlockKey) < 32 && len(s.BlockKey)%8 != 0 {
		return ErrInvalidSessionBlockKey
	}

	// HMAC wants 32 or 64, we're going to err on the side of caution and force 64.
	if len(s.HashKey) != 64 {
		return ErrInvalidSessionHashKey
	}

	if s.Name == "" {
		s.Name = "session"
	}

	if s.Path == "" {
		s.Path = "/"
	}

	if s.MaxAge == 0 {
		s.MaxAge = 60 * 60 // 1 hour
	}

	if s.Domain == "" {
		return ErrMissingSessionDomain
	}

	return nil
}

func IsValidOAuth2Provider(provider string) bool {
	provider = strings.ToLower(provider)
	return provider == "vatsim" || provider == "adh"
}
