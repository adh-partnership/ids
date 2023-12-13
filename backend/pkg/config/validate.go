package config

import (
	"errors"
)

var (
	ErrInvalidCacheDriver    = errors.New("invalid cache driver")
	ErrInvalidCacheRedisHost = errors.New("invalid redis host or port")
	ErrInvalidCacheRedisDB   = errors.New("invalid redis db")
)

func ValidateConfig(c *Config) error {
	if err := ValidateCache(&c.Cache); err != nil {
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
