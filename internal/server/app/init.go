//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	commonconfig "github.com/ognick/word_of_wisdom/internal/common/config"
	internalconfig "github.com/ognick/word_of_wisdom/internal/server/internal/config"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/challenge"
	challengeusecase "github.com/ognick/word_of_wisdom/internal/server/internal/services/challenge/usecase"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom"
	wisdomtcpV1 "github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/api/tcp/v1"
	"github.com/ognick/word_of_wisdom/pkg/http"
	"github.com/ognick/word_of_wisdom/pkg/logger"
	"github.com/ognick/word_of_wisdom/pkg/logger/zap"
	"github.com/ognick/word_of_wisdom/pkg/pow"
	"github.com/ognick/word_of_wisdom/pkg/tcp"
)

func InitializeApp() (*App, error) {
	wire.Build(Application)
	return nil, nil
}

func provideLogger(cfg commonconfig.Config) logger.Logger {
	return zap.NewLogger(cfg.Logger)
}

func provideChallengeTimeout(cfg commonconfig.Config) wisdomtcpV1.ChallengeTimeout {
	return wisdomtcpV1.ChallengeTimeout(cfg.ChallengeTimeout)
}

func provideProofOfWorkGenerator(cfg internalconfig.Config) challengeusecase.ProofOfWorkGenerator {
	return pow.NewGenerator(cfg.ChallengeComplexity)
}

func provideTCPAddr(cfg commonconfig.Config) tcp.Addr {
	return tcp.Addr(cfg.TCPAddress)
}

func provideHTTPAddr(cfg commonconfig.Config) http.Addr {
	return http.Addr(cfg.HTTPAddress)
}

var Application = wire.NewSet(
	// Application
	NewApp,
	commonconfig.NewConfig,
	provideLogger,
	tcp.NewServer,
	provideTCPAddr,
	http.NewServer,
	provideHTTPAddr,
	registerHTTPHandlers,
	// Domain
	wisdom.Init,
	challenge.Init,
	provideChallengeTimeout,
	provideProofOfWorkGenerator,
	internalconfig.NewConfig,
)
