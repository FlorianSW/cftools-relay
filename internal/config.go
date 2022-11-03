package internal

import (
	"cftools-relay/internal/domain"
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
)

type Discord struct {
	WebhookUrl string `json:"webhook_url"`
}

type History struct {
	StoragePath string `json:"storage_path"`
}

type Config struct {
	Port    int                      `json:"port"`
	Secret  string                   `json:"secret,omitempty"`
	Servers map[string]domain.Server `json:"servers"`
	Discord Discord                  `json:"discord"`
	History History                  `json:"history"`
	Filter  domain.FilterList        `json:"filter"`
}

func NewConfig(path string, logger lager.Logger) (Config, error) {
	config, err := readConfig(path, logger)
	if err != nil {
		return config, err
	}

	if config.Filter == nil {
		config.Filter = domain.FilterList{}
	} else {
		for i, filter := range config.Filter {
			if (filter.Color != "" || filter.Message != "") && (filter.Format == nil || filter.Format.Type == "") {
				config.Filter[i].Format = &domain.Format{
					Type: domain.FormatTypeRich,
					Parameters: map[string]interface{}{
						"color":   filter.Color,
						"message": filter.Message,
					},
				}
				config.Filter[i].Color = ""
				config.Filter[i].Message = ""
			}
		}
	}
	if config.History.StoragePath == "" {
		config.History.StoragePath = "./storage"
	}
	if len(config.Servers) != 0 && config.Secret != "" {
		return config, errors.New("can not have a secret and servers configured at the same time")
	}
	if config.Secret != "" {
		config.Servers = map[string]domain.Server{}
		config.Servers[""] = domain.Server{Secret: config.Secret}
		config.Secret = ""
	}
	for name, _ := range config.Servers {
		if url.PathEscape(name) != name {
			return config, fmt.Errorf("%s is expected to be URL-safe", name)
		}
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
