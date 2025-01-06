//go:build wireinject
// +build wireinject

package app

import (
	"fmt"

	"github.com/google/wire"
	commonconfig "github.com/ognick/word_of_wisdom/internal/common/config"
	privateconfig "github.com/ognick/word_of_wisdom/internal/server/internal/config"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/challenge"
	challengeusecase "github.com/ognick/word_of_wisdom/internal/server/internal/services/challenge/usecase"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom"
	wisdomtcpV1 "github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/api/tcp/v1"
	"github.com/ognick/word_of_wisdom/pkg/http"
	"github.com/ognick/word_of_wisdom/pkg/logger"
	"github.com/ognick/word_of_wisdom/pkg/pow"
	"github.com/ognick/word_of_wisdom/pkg/tcp"
)

func InitializeApp() (*App, error) {
	wire.Build(Application)
	return nil, nil
}

func provideLogger(cfg commonconfig.Config) (logger.Logger, error) {
	l := logger.NewLogger()
	if err := logger.SetLogLevel(cfg.LogLevel); err != nil {
		return l, fmt.Errorf("failed set log level: %v", err)
	}
	return l, nil
}

func provideChallengeComplexity(cfg privateconfig.Config) pow.Complexity {
	return cfg.ChallengeComplexity
}

func provideChallengeTimeout(cfg commonconfig.Config) wisdomtcpV1.ChallengeTimeout {
	return wisdomtcpV1.ChallengeTimeout(cfg.ChallengeTimeout)
}

func provideProofOfWorkGenerator(complexity pow.Complexity) challengeusecase.ProofOfWorkGenerator {
	return pow.NewGenerator(complexity)
}

func provideTCPAddress(cfg commonconfig.Config) tcp.Address {
	return tcp.Address(cfg.TCPAddress)
}

func provideHTTPAddr(cfg commonconfig.Config) http.Address {
	return http.Address(cfg.HTTPAddress)
}

var Application = wire.NewSet(
	NewApp,

	provideLogger,

	commonconfig.Set,
	privateconfig.Set,

	wisdom.Set,
	challenge.Set,
	provideChallengeComplexity,
	provideChallengeTimeout,
	provideProofOfWorkGenerator,

	tcp.NewServer,
	provideTCPAddress,

	http.NewServer,
	provideHTTPAddr,
)
