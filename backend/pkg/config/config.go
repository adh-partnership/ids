package config

import (
	"encoding/json"
	"os"
	"sync"

	"sigs.k8s.io/yaml"
)

var (
	cfg      *Config
	airports []string
	mutex    sync.RWMutex
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

	ValidateConfig(c)

	cfg = c

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
