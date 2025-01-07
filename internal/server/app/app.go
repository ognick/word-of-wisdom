package app

import (
	"log"

	"github.com/ognick/word_of_wisdom/internal/server/internal/config"
	"github.com/ognick/word_of_wisdom/pkg/logger"
	"github.com/ognick/word_of_wisdom/pkg/shutdown"
)

type App struct {
	log logger.Logger
	cfg config.Config
}

func NewApp(log logger.Logger, cfg config.Config) *App {
	return &App{
		log: log,
		cfg: cfg,
	}
}

func (app *App) Run() {
	tcpSrv, err := initTCPServer()
	if err != nil {
		app.log.Fatalf("failed to initialize TCP server: %v", err)
	}

	httpSrv, err := initHttpServer()
	if err != nil {
		app.log.Fatalf("failed to initialize HTTP server: %v", err)
	}

	runner, gracefulCtx := shutdown.CreateRunnerWithGracefulContext()
	runner.Go(func() error {
		return tcpSrv.Run(gracefulCtx)
	})

	runner.Go(func() error {
		return httpSrv.Run(gracefulCtx)
	})

	app.log.Infof("Server started")

	// Awaiting graceful shutdown
	if err := runner.Wait(); err != nil {
		log.Fatalf("%v", err)
	}

	app.log.Infof("Server gracefully finished")
}
