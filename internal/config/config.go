package config

import (
	"encoding/json"
	"github.com/ent1k1377/load_balancer/internal/logger"
	"os"
)

type Config struct {
	Port                int      `json:"port"`
	Backends            []string `json:"backends"`
	DefaultCapacity     int      `json:"defaultCapacity"`
	DefaultRefillRate   int      `json:"defaultRefillRate"`
	DefaultRefillPeriod int      `json:"defaultRefillPeriodMs"`
	Strategy            string   `json:"strategy"`
	HealthCheckInterval int      `json:"healthCheckIntervalSec"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Errorf("failed to read config file: %v", err)
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		logger.Errorf("failed to unmarshal config file: %v", err)
		return nil, err
	}

	logger.Info("config loaded successfully")
	return &config, nil
}
