package app

import (
	"context"
	"errors"
	"time"

	"word_of_wisdom/internal/client/api/tcp/v1"
	"word_of_wisdom/internal/client/service"
	"word_of_wisdom/internal/config"
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
		log.Fatalf("failed to set log level: %v", err)
	}

	// Services
	solverService := service.NewSolverService()

	// Handler
	handler := v1.NewHandler(solverService)

	// TCP Client
	client := tcp.NewClient(cfg.HTTPAddress, handler.Handle)

	// Run
	runner, ctx := shutdown.CreateRunnerWithGracefulContext()
	runner.Go(func() error {
		for {
			if err := client.Connect(ctx); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				log.Errorf("failed to connect: %v", err)
				<-time.After(cfg.ChallengeTimeout)
			}
		}
	})

	// Awaiting graceful shutdown
	if err := runner.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("%v", err)
	}

	log.Infof("Client gracefully finished")
}
