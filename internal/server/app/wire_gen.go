// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/ognick/word_of_wisdom/internal/common/config"
	"github.com/ognick/word_of_wisdom/internal/server/docs"
	config2 "github.com/ognick/word_of_wisdom/internal/server/internal/config"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/challenge"
	challenge2 "github.com/ognick/word_of_wisdom/internal/server/internal/services/challenge/usecase"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/api/http/v1"
	v1_2 "github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/api/tcp/v1"
	"github.com/ognick/word_of_wisdom/pkg/http"
	"github.com/ognick/word_of_wisdom/pkg/logger"
	"github.com/ognick/word_of_wisdom/pkg/logger/zap"
	"github.com/ognick/word_of_wisdom/pkg/pow"
	"github.com/ognick/word_of_wisdom/pkg/tcp"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	http2 "net/http"
)

// Injectors from init.go:

func InitializeApp() (*App, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	logger := provideLogger(configConfig)
	config3, err := config2.NewConfig()
	if err != nil {
		return nil, err
	}
	app := NewApp(logger, config3)
	return app, nil
}

// Injectors from inject_http.go:

func initHttpServer() (*http.Server, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	logger := provideLogger(configConfig)
	addr := provideHTTPAddr(configConfig)
	config3, err := config2.NewConfig()
	if err != nil {
		return nil, err
	}
	proofOfWorkGenerator := provideProofOfWorkGenerator(config3)
	usecasesChallenge := challenge.ProvideUsecase(proofOfWorkGenerator)
	repositoriesWisdom := wisdom.ProvideRepo()
	usecasesWisdom := wisdom.ProvideUsecase(repositoriesWisdom)
	handler := v1.NewHandler(logger, usecasesChallenge, usecasesWisdom)
	httpHandler := registerHTTPHandlers(handler)
	server := http.NewServer(logger, addr, httpHandler)
	return server, nil
}

// Injectors from inject_tcp.go:

func initTCPServer() (*tcp.Server, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	logger := provideLogger(configConfig)
	addr := provideTCPAddr(configConfig)
	config3, err := config2.NewConfig()
	if err != nil {
		return nil, err
	}
	proofOfWorkGenerator := provideProofOfWorkGenerator(config3)
	usecasesChallenge := challenge.ProvideUsecase(proofOfWorkGenerator)
	repositoriesWisdom := wisdom.ProvideRepo()
	usecasesWisdom := wisdom.ProvideUsecase(repositoriesWisdom)
	challengeTimeout := provideChallengeTimeout(configConfig)
	handler := v1_2.NewHandler(logger, usecasesChallenge, usecasesWisdom, challengeTimeout)
	v := v1_2.ProvideTCPHandle(handler)
	server := tcp.NewServer(logger, addr, v)
	return server, nil
}

// init.go:

func provideLogger(cfg config.Config) logger.Logger {
	return zap.NewLogger(cfg.Logger)
}

func provideChallengeTimeout(cfg config.Config) v1_2.ChallengeTimeout {
	return v1_2.ChallengeTimeout(cfg.ChallengeTimeout)
}

func provideProofOfWorkGenerator(cfg config2.Config) challenge2.ProofOfWorkGenerator {
	return pow.NewGenerator(cfg.ChallengeComplexity)
}

var Application = wire.NewSet(

	NewApp, config.NewConfig, provideLogger, wisdom.Init, challenge.Init, provideChallengeTimeout,
	provideProofOfWorkGenerator, config2.NewConfig,
)

// inject_http.go:

// @title Word of Wisdom API
// @version 1.0
// @description This is a simple API for getting wisdoms
// @contact.name Dmitry Aleksandrov
// @contact.email ogneslav.work@gmail.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
func registerHTTPHandlers(
	wisdomV1 *v1.Handler,
) http2.Handler {
	router := gin.Default()
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.String(http2.StatusOK, "ok")
	})

	wisdomV1.Register(router)
	docs.SwaggerInfo.
		BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}

func provideHTTPAddr(cfg config.Config) http.Addr {
	return http.Addr(cfg.HTTPAddress)
}

// inject_tcp.go:

func provideTCPAddr(cfg config.Config) tcp.Addr {
	return tcp.Addr(cfg.TCPAddress)
}
