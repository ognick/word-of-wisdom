package app

import (
	commonconfig "word_of_wisdom/internal/common/config"
	httpV1 "word_of_wisdom/internal/server/internal/api/http/v1"
	tcpV1 "word_of_wisdom/internal/server/internal/api/tcp/v1"
	privateconfig "word_of_wisdom/internal/server/internal/config"
	"word_of_wisdom/internal/server/internal/repository"
	"word_of_wisdom/internal/server/internal/service"
	"word_of_wisdom/pkg/http"
	"word_of_wisdom/pkg/logger"
	"word_of_wisdom/pkg/shutdown"
	"word_of_wisdom/pkg/tcp"
)

func Run() {
	log := logger.NewLogger()
	commonCfg, err := commonconfig.NewConfig()
	if err != nil {
		log.Fatalf("failed to init common config: %v", err)
	}
	if err := logger.SetLogLevel(commonCfg.LogLevel); err != nil {
		log.Fatalf("failed set log level: %v", err)
	}

	privateCfg, err := privateconfig.NewConfig()
	if err != nil {
		log.Fatalf("failed to init private config: %v", err)
	}

	//Repositories
	wisdomRepo := repository.NewWisdomRepository()

	// Services
	challengeService := service.NewChallengeService(privateCfg.ChallengeComplexity)
	wisdomService := service.NewWisdomService(wisdomRepo)

	// TCP Handlers
	tcpHandler := tcpV1.NewHandler(
		challengeService,
		wisdomService,
		commonCfg.ChallengeTimeout,
	)

	// TCP Server
	tcpSrv := tcp.NewServer(commonCfg.TCPAddress, tcpHandler.Handle)

	// HTTP Handlers
	httpHandler := httpV1.NewHandler(
		challengeService,
		wisdomService,
	).Init()

	// HTTP Server
	httpSrv := http.NewServer(commonCfg.HTTPAddress, httpHandler)

	// Run
	runner, gracefulCtx := shutdown.CreateRunnerWithGracefulContext()
	runner.Go(func() error {
		return tcpSrv.Run(gracefulCtx)
	})
	runner.Go(func() error {
		return httpSrv.Run(gracefulCtx)
	})

	log.Infof("Server started")

	// Awaiting graceful shutdown
	if err := runner.Wait(); err != nil {
		log.Fatalf("%v", err)
	}

	log.Infof("Server gracefully finished")
}
