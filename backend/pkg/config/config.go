package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/adh-partnership/api/pkg/logger"
	"sigs.k8s.io/yaml"
)

var (
	cfg      *Config
	airports []string
	mutex    sync.RWMutex
	log      = logger.Logger.WithField("component", "config")
)

func ParseAirports(file string) error {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &airports)
	if err != nil {
		return err
	}

	return nil
}

func ParseConfig(file string) error {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	c := &Config{}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	cfg = c
	jd, _ := json.Marshal(cfg)
	log.Debugf("Configuration: %s", jd)

	return nil
}

func GetAirports() []string {
	mutex.RLock()
	defer mutex.RUnlock()

	return airports
}

func GetConfig() *Config {
	mutex.RLock()
	defer mutex.RUnlock()

	return cfg
}
