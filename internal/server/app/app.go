package app

import (
	"log"

	"github.com/ognick/word_of_wisdom/pkg/http"
	"github.com/ognick/word_of_wisdom/pkg/logger"
	"github.com/ognick/word_of_wisdom/pkg/shutdown"
	"github.com/ognick/word_of_wisdom/pkg/tcp"
)

type App struct {
	log        logger.Logger
	tcpServer  *tcp.Server
	httpServer *http.Server
}

func NewApp(
	log logger.Logger,
	tcpServer *tcp.Server,
	httpServer *http.Server,
) *App {
	return &App{
		log:        log,
		tcpServer:  tcpServer,
		httpServer: httpServer,
	}
}

func (app *App) Run() {
	runner, gracefulCtx := shutdown.CreateRunnerWithGracefulContext()
	runner.Go(func() error {
		return app.tcpServer.Run(gracefulCtx)
	})

	runner.Go(func() error {
		return app.httpServer.Run(gracefulCtx)
	})

	app.log.Infof("Server started")

	// Awaiting graceful shutdown
	if err := runner.Wait(); err != nil {
		log.Fatalf("%v", err)
	}

	app.log.Infof("Server gracefully finished")
}
