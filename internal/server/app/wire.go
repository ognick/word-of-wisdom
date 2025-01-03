//go:build wireinject
// +build wireinject

package app

import (
	"fmt"
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	commonconfig "word_of_wisdom/internal/common/config"
	privateconfig "word_of_wisdom/internal/server/internal/config"
	"word_of_wisdom/internal/server/internal/domain/types/repositories"
	"word_of_wisdom/internal/server/internal/domain/types/usecases"
	challengeusecase "word_of_wisdom/internal/server/internal/services/challenge/usecase"
	wisdomhttpV1 "word_of_wisdom/internal/server/internal/services/wisdom/api/http/v1"
	wisdomtcpV1 "word_of_wisdom/internal/server/internal/services/wisdom/api/tcp/v1"
	wisdomrepo "word_of_wisdom/internal/server/internal/services/wisdom/repository"
	wisdomusecase "word_of_wisdom/internal/server/internal/services/wisdom/usecase"
	"word_of_wisdom/pkg/http"
	"word_of_wisdom/pkg/logger"
	"word_of_wisdom/pkg/pow"
	"word_of_wisdom/pkg/tcp"
)

func InitApp() (*App, error) {
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

	// TCP Handlers
	tcp.NewServer,
	wisdomtcpV1.Set,
	provideTCPAddress,
	provideChallengeTimeout,

	// HTTP Handlers
	http.NewServer,
	wisdomhttpV1.Set, wire.Bind(new(nethttp.Handler), new(*gin.Engine)),
	provideHTTPAddr,

	// Repositories
	wisdomrepo.NewRepository, wire.Bind(new(repositories.Wisdom), new(*wisdomrepo.Repository)),
	pow.NewGenerator, wire.Bind(new(challengeusecase.ProofOfWorkGenerator), new(*pow.Generator)),
	provideChallengeComplexity,

	// Usecases
	challengeusecase.NewUsecase, wire.Bind(new(usecases.Challenge), new(*challengeusecase.Usecase)),
	wisdomusecase.NewUsecase, wire.Bind(new(usecases.Wisdom), new(*wisdomusecase.Usecase)),
)
