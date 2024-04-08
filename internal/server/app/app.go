package app

import (
	"word_of_wisdom/internal/config"
	"word_of_wisdom/internal/server/api/tcp/v1"
	"word_of_wisdom/internal/server/repository"
	"word_of_wisdom/internal/server/service"
	"word_of_wisdom/pkg/logger"
	"word_of_wisdom/pkg/shutdown"
	"word_of_wisdom/pkg/tcp"
)

func Run() {
	log := logger.NewLogger()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}
	if err := logger.SetLogLevel(cfg.LogLevel); err != nil {
		log.Fatalf("failed set log level: %v", err)
	}

	//Repositories
	wisdomRepo := repository.NewWisdomRepository()

	// Services
	challengeService := service.NewChallengeService(cfg.ChallengeComplexity)
	wisdomService := service.NewWisdomService(wisdomRepo)

	// Handlers
	handler := v1.NewHandler(
		challengeService,
		wisdomService,
		cfg.ChallengeTimeout,
	)

	// TCP Server
	srv := tcp.NewServer(cfg.HTTPAddress, handler.Handle)

	// Run
	runner, gracefulCtx := shutdown.CreateRunnerWithGracefulContext()
	runner.Go(func() error {
		return srv.Run(gracefulCtx)
	})

	log.Infof("Server started")

	// Awaiting graceful shutdown
	if err := runner.Wait(); err != nil {
		log.Fatalf("%v", err)
	}

	log.Infof("Server gracefully finished")
}
