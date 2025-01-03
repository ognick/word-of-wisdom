// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	http2 "net/http"
	"word_of_wisdom/internal/common/config"
	config2 "word_of_wisdom/internal/server/internal/config"
	"word_of_wisdom/internal/server/internal/domain/types/repositories"
	"word_of_wisdom/internal/server/internal/domain/types/usecases"
	"word_of_wisdom/internal/server/internal/services/challenge/usecase"
	v1_2 "word_of_wisdom/internal/server/internal/services/wisdom/api/http/v1"
	"word_of_wisdom/internal/server/internal/services/wisdom/api/tcp/v1"
	"word_of_wisdom/internal/server/internal/services/wisdom/repository"
	"word_of_wisdom/internal/server/internal/services/wisdom/usecase"
	"word_of_wisdom/pkg/http"
	"word_of_wisdom/pkg/logger"
	"word_of_wisdom/pkg/pow"
	"word_of_wisdom/pkg/tcp"
)

// Injectors from wire.go:

func InitializeAppWithWire() (*App, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	logger, err := provideLogger(configConfig)
	if err != nil {
		return nil, err
	}
	address := provideTCPAddress(configConfig)
	config3, err := config2.NewConfig()
	if err != nil {
		return nil, err
	}
	complexity := provideChallengeComplexity(config3)
	generator := pow.NewGenerator(complexity)
	usecase := challenge.NewUsecase(generator)
	repositoryRepository := repository.NewRepository()
	wisdomUsecase := wisdom.NewUsecase(repositoryRepository)
	challengeTimeout := provideChallengeTimeout(configConfig)
	handler := v1.NewHandler(usecase, wisdomUsecase, challengeTimeout)
	v := v1.ProvideTCPHandle(handler)
	server := tcp.NewServer(address, v)
	httpAddress := provideHTTPAddr(configConfig)
	v1Handler := v1_2.NewHandler(usecase, wisdomUsecase)
	engine := v1_2.ProvideHTTPEngine(v1Handler)
	httpServer := http.NewServer(httpAddress, engine)
	app := NewApp(logger, server, httpServer)
	return app, nil
}

// wire.go:

func provideLogger(cfg config.Config) (logger.Logger, error) {
	l := logger.NewLogger()
	if err := logger.SetLogLevel(cfg.LogLevel); err != nil {
		return l, fmt.Errorf("failed set log level: %v", err)
	}
	return l, nil
}

func provideChallengeComplexity(cfg config2.Config) pow.Complexity {
	return cfg.ChallengeComplexity
}

func provideChallengeTimeout(cfg config.Config) v1.ChallengeTimeout {
	return v1.ChallengeTimeout(cfg.ChallengeTimeout)
}

func provideTCPAddress(cfg config.Config) tcp.Address {
	return tcp.Address(cfg.TCPAddress)
}

func provideHTTPAddr(cfg config.Config) http.Address {
	return http.Address(cfg.HTTPAddress)
}

var Application = wire.NewSet(
	NewApp,

	provideLogger, config.Set, config2.Set, tcp.NewServer, v1.Set, provideTCPAddress,
	provideChallengeTimeout, http.NewServer, v1_2.Set, wire.Bind(new(http2.Handler), new(*gin.Engine)), provideHTTPAddr, repository.NewRepository, wire.Bind(new(repositories.Wisdom), new(*repository.Repository)), pow.NewGenerator, wire.Bind(new(challenge.ProofOfWorkGenerator), new(*pow.Generator)), provideChallengeComplexity, challenge.NewUsecase, wire.Bind(new(usecases.Challenge), new(*challenge.Usecase)), wisdom.NewUsecase, wire.Bind(new(usecases.Wisdom), new(*wisdom.Usecase)),
)
