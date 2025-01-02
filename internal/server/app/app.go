package app

import (
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

	// Proof of concept generator
	proofOfWorkGenerator := pow.NewGenerator(privateCfg.ChallengeComplexity)

	// Repositories
	wisdomRepo := wisdomrepo.NewWisdomRepository()

	// Services
	challengeService := challengeusecase.NewUsecase(proofOfWorkGenerator)
	wisdomService := wisdomusecase.NewService(wisdomusecase.DepRepos{
		Wisdom: wisdomRepo,
	})

	// TCP Handlers
	tcpHandler := wisdomtcpV1.NewHandler(
		challengeService,
		wisdomService,
		commonCfg.ChallengeTimeout,
	)

	// TCP Server
	tcpSrv := tcp.NewServer(commonCfg.TCPAddress, tcpHandler.Handle)

	// HTTP Handlers
	httpHandler := wisdomhttpV1.NewHandler(
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
