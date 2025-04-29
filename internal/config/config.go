package config

import (
	"encoding/json"
	"github.com/ent1k1377/load_balancer/internal/logger"
	"os"
)

// Config содержит параметры конфигурации для балансировщика нагрузки.
type Config struct {
	Port                int      `json:"port"`                   // Порт, на котором будет работать балансировщик нагрузки
	Backends            []string `json:"backends"`               // Список URL backend-серверов
	DefaultCapacity     int      `json:"defaultCapacity"`        // Ёмкость лимитера запросов по умолчанию
	DefaultRefillRate   int      `json:"defaultRefillRate"`      // Скорость пополнения лимита по умолчанию
	DefaultRefillPeriod int      `json:"defaultRefillPeriodMs"`  // Период пополнения лимита (в миллисекундах)
	Strategy            string   `json:"strategy"`               // Стратегия балансировки нагрузки (например, "round_robin" или "least_connections")
	HealthCheckInterval int      `json:"healthCheckIntervalSec"` // Интервал для проверки состояния backend'ов (в секундах)
}

// LoadConfig загружает конфигурацию из файла по указанному пути и возвращает объект конфигурации.

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
