package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/google/wire"

	"github.com/ognick/word_of_wisdom/pkg/http"
	"github.com/ognick/word_of_wisdom/pkg/tcp"
)

type Config struct {
	TCPAddress       tcp.Address   `env:"CONF_TCP_ADDRESS" envDefault:":8080"`
	HTTPAddress      http.Address  `env:"CONF_HTTP_ADDRESS" envDefault:":8888"`
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
