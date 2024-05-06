package config

import (
	"fmt"

	"github.com/caarlos0/env/v7"
)

type Config struct {
	ChallengeComplexity byte `env:"CONF_CHALLENGE_COMPLEXITY" envDefault:"4"`
}

func NewConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse env: %w", err)
	}

	return cfg, nil
}
