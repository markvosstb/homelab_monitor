// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"time"

	yaml "go.yaml.in/yaml/v4"
)

type ServiceConfig struct {
	Name string `yaml:"name"`
	Subdomain string `yaml:"subdomain"`
}

type Config struct {
	Interval time.Duration `yaml:"interval"`
	Name string `yaml:"name"`
	Site string `yaml:"site"`
	Services []ServiceConfig `yaml:"services"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	cfg := Config{}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Services) == 0 {
		return fmt.Errorf("at least one service must be configured")
	}

	for _, s := range c.Services {
		if s.Name == "" || s.Subdomain == "" {
			return fmt.Errorf("services entry missing required field (name, subdomain)")
		}
	}
	return nil
}
