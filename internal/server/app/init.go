//go:build wireinject
// +build wireinject

package app

import (
	"fmt"
	"net"
	nethttp "net/http"

	"github.com/google/wire"

	commonconfig "github.com/ognick/word_of_wisdom/internal/common/config"
	internalconfig "github.com/ognick/word_of_wisdom/internal/server/internal/config"
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

func provideChallengeTimeout(cfg commonconfig.Config) wisdomtcpV1.ChallengeTimeout {
	return wisdomtcpV1.ChallengeTimeout(cfg.ChallengeTimeout)
}

func provideProofOfWorkGenerator(cfg internalconfig.Config) challengeusecase.ProofOfWorkGenerator {
	return pow.NewGenerator(cfg.ChallengeComplexity)
}

func provideTCPServer(cfg commonconfig.Config, handler func(conn net.Conn)) *tcp.Server {
	return tcp.NewServer(cfg.TCPAddress, handler)
}

func provideHTTPServer(cfg commonconfig.Config, handler nethttp.Handler) *http.Server {
	return http.NewServer(cfg.HTTPAddress, handler)
}

var Application = wire.NewSet(
	// Application
	NewApp,
	commonconfig.NewConfig,
	provideLogger,
	provideTCPServer,
	provideHTTPServer,
	registerHTTPHandlers,
	// Domain
	wisdom.Init,
	challenge.Init,
	provideChallengeTimeout,
	provideProofOfWorkGenerator,
	internalconfig.NewConfig,
)
