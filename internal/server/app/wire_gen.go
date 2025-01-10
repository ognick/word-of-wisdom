// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/ognick/word_of_wisdom/internal/common/config"
	config2 "github.com/ognick/word_of_wisdom/internal/server/internal/config"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/challenge"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/api/http/v1"
	v1_2 "github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/api/tcp/v1"
	"github.com/ognick/word_of_wisdom/pkg/http"
	"github.com/ognick/word_of_wisdom/pkg/lifecycle"
	"github.com/ognick/word_of_wisdom/pkg/tcp"
)

import (
	_ "github.com/ognick/word_of_wisdom/internal/server/docs"
)

// Injectors from inject_app.go:

func InitializeApp() (*App, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	logger := provideLogger(configConfig)
	lifecycleLifecycle := lifecycle.NewLifecycle()
	addr := provideHTTPAddr(configConfig)
	config3, err := config2.NewConfig()
	if err != nil {
		return nil, err
	}
	proofOfWorkGenerator := provideProofOfWorkGenerator(config3)
	usecasesChallenge := challenge.ProvideUsecase(proofOfWorkGenerator)
	repo := wisdom.ProvideRepo()
	usecasesWisdom := wisdom.ProvideUsecase(repo)
	handler := v1.NewHandler(logger, usecasesChallenge, usecasesWisdom)
	httpHandler := registerHTTPHandlers(handler)
	server := http.NewServer(lifecycleLifecycle, logger, addr, httpHandler)
	tcpAddr := provideTCPAddr(configConfig)
	challengeTimeout := provideChallengeTimeout(configConfig)
	v1Handler := v1_2.NewHandler(logger, usecasesChallenge, usecasesWisdom, challengeTimeout)
	v := v1_2.ProvideTCPHandle(v1Handler)
	tcpServer := tcp.NewServer(lifecycleLifecycle, logger, tcpAddr, v)
	modules := Modules{
		httpServer: server,
		tcpServer:  tcpServer,
	}
	app := NewApp(logger, lifecycleLifecycle, modules)
	return app, nil
}
