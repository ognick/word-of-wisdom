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

	tcp.NewServer,
	wisdomtcpV1.Set,
	provideTCPAddress,
	provideChallengeTimeout,
	http.NewServer,
	wisdomhttpV1.Set,
	provideHTTPAddr,

	commonconfig.Set,
	privateconfig.Set,

	provideLogger,

	wisdomrepo.Set,
	pow.NewGenerator,
	provideChallengeComplexity,

	challengeusecase.Set,
	wisdomusecase.Set,

	wire.Bind(new(wisdomtcpV1.ChallengeUsecase), new(*challengeusecase.Usecase)),
	wire.Bind(new(wisdomhttpV1.ChallengeUsecase), new(*challengeusecase.Usecase)),
	wire.Bind(new(wisdomtcpV1.WisdomUsecase), new(*wisdomusecase.Usecase)),
	wire.Bind(new(wisdomhttpV1.WisdomUsecase), new(*wisdomusecase.Usecase)),
	wire.Bind(new(wisdomusecase.WisdomRepo), new(*wisdomrepo.WisdomRepository)),
	wire.Bind(new(challengeusecase.ProofOfWorkGenerator), new(*pow.Generator)),
	wire.Bind(new(nethttp.Handler), new(*gin.Engine)),
)
