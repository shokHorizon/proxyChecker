package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type (
	Config struct {
		Sources []Source `yaml:"sources"`
		Pinger  pinger   `yaml:"pinger"`
		Log     log      `yaml:"log"`
	}

	Source struct {
		Regexp  string   `env-required:"true" yaml:"regexp" env:"SOURCE_REGEXP"`
		Headers []string `env-required:"false" yaml:"headers" env:"SOURCE_HEADERS"`
		Urls    []string `env-required:"true" yaml:"urls" env:"SOURCE_URLS"`
		Name    string   `env-required:"true" yaml:"name" env:"SOURCE_NAME"`
	}

	pinger struct {
		Timeout time.Duration `env-required:"false" yaml:"timeout" env:"PINGER_TIMEOUT"`
		Workers int           `env-required:"false" yaml:"workers" env:"PINGER_WORKERS"`
	}

	log struct {
		Level string `env-required:"true" yaml:"level"   env:"LOG_LEVEL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
