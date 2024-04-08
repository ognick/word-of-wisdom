package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v7"
)

type Config struct {
	HTTPAddress         string        `env:"CONF_HTTP_ADDRESS" envDefault:":8080"`
	LogLevel            string        `env:"CONF_LOG_LEVEL" envDefault:"debug"`
	ChallengeTimeout    time.Duration `env:"CONF_TIMEOUT" envDefault:"1s"`
	ChallengeComplexity byte          `env:"CONF_CHALLENGE_COMPLEXITY" envDefault:"4"`
}

func NewConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse env: %w", err)
	}

	return cfg, nil
}
