package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v7"
)

type Logger struct {
	Development   bool   `env:"DEVELOPMENT" envDefault:"false"`
	DisableCaller bool   `env:"DISABLE_CALLER" envDefault:"true"`
	DisableJson   bool   `env:"ENCODING" envDefault:"true"`
	Level         string `env:"LEVEL" envDefault:"debug"`
}

type Config struct {
	TCPAddress       string        `env:"CONF_TCP_ADDRESS" envDefault:":8080"`
	HTTPAddress      string        `env:"CONF_HTTP_ADDRESS" envDefault:":8888"`
	Logger           Logger        `env:"CONF_LOGGER"`
	ChallengeTimeout time.Duration `env:"CONF_TIMEOUT" envDefault:"1s"`
}

func NewConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse env: %w", err)
	}

	return cfg, nil
}
