package internal

import (
	"cftools-relay/internal/domain"
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"os"
)

type Discord struct {
	WebhookUrl string `json:"webhook_url"`
}

type Config struct {
	Port    int               `json:"port"`
	Secret  string            `json:"secret"`
	Discord Discord           `json:"discord"`
	Filter  domain.FilterList `json:"filter"`
}

func NewConfig(path string, logger lager.Logger) (Config, error) {
	config, err := readConfig(path, logger)
	if err != nil {
		return config, err
	}

	if config.Filter == nil {
		config.Filter = domain.FilterList{}
	}

	return config, persistConfig(path, config)
}

func readConfig(path string, logger lager.Logger) (Config, error) {
	var config Config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Info("create-config")
		config = Config{
			Port: 8080,
		}
	} else {
		logger.Info("read-existing-config")
		c, err := os.ReadFile(path)
		if err != nil {
			return Config{}, err
		}
		err = json.Unmarshal(c, &config)
		if err != nil {
			return Config{}, err
		}
	}
	return config, nil
}

func persistConfig(path string, config Config) error {
	c, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, c, 0655)
}
