package app

import (
	"github.com/google/wire"

	commonconfig "github.com/ognick/word_of_wisdom/internal/common/config"
	internalconfig "github.com/ognick/word_of_wisdom/internal/server/internal/config"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/challenge"
	challengeusecase "github.com/ognick/word_of_wisdom/internal/server/internal/services/challenge/usecase"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom"
	wisdomtcpV1 "github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/api/tcp/v1"
	"github.com/ognick/word_of_wisdom/pkg/lifecycle"
	"github.com/ognick/word_of_wisdom/pkg/logger"
	"github.com/ognick/word_of_wisdom/pkg/logger/zap"
	"github.com/ognick/word_of_wisdom/pkg/pow"
)

func provideLogger(cfg commonconfig.Config) logger.Logger {
	return zap.NewLogger(zap.Config{
		Level:         cfg.Logger.Level,
		Development:   cfg.Logger.Development,
		DisableCaller: cfg.Logger.DisableCaller,
		DisableJson:   cfg.Logger.DisableJson,
	})
}

func provideChallengeTimeout(cfg commonconfig.Config) wisdomtcpV1.ChallengeTimeout {
	return wisdomtcpV1.ChallengeTimeout(cfg.ChallengeTimeout)
}

func provideProofOfWorkGenerator(cfg internalconfig.Config) challengeusecase.ProofOfWorkGenerator {
	return pow.NewGenerator(cfg.ChallengeComplexity)
}

var Application = wire.NewSet(
	// Application
	NewApp,
	lifecycle.NewLifecycle,
	commonconfig.NewConfig,
	provideLogger,
	// Domain
	wisdom.Init,
	challenge.Init,
	provideChallengeTimeout,
	provideProofOfWorkGenerator,
	internalconfig.NewConfig,
	// Runnable modules
	wire.Struct(new(runnableModules), "*"),
	initHTTPServer,
	initTCPServer,
)
