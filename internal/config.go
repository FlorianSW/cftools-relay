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

	if config.Filter == nil {
		config.Filter = domain.FilterList{}
	}
	c, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return Config{}, err
	}
	err = os.WriteFile(path, c, 0655)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
