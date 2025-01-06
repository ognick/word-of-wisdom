package config

import (
	"fmt"

	"github.com/caarlos0/env/v7"
	"github.com/google/wire"

	"github.com/ognick/word_of_wisdom/pkg/pow"
)

type Config struct {
	ChallengeComplexity pow.Complexity `env:"CONF_CHALLENGE_COMPLEXITY" envDefault:"4"`
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
