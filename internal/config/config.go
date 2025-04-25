package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port     string   `json:"port"`
	Backends []string `json:"backends"`
}

func LoadConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	return &config
}
