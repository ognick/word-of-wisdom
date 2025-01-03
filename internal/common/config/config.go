package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/google/wire"
)

type Config struct {
	TCPAddress       string        `env:"CONF_TCP_ADDRESS" envDefault:":8080"`
	HTTPAddress      string        `env:"CONF_HTTP_ADDRESS" envDefault:":8888"`
	LogLevel         string        `env:"CONF_LOG_LEVEL" envDefault:"debug"`
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

var Set = wire.NewSet(NewConfig)
