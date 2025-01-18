package app

import (
	"log"

	"github.com/ognick/word_of_wisdom/pkg/http"
	"github.com/ognick/word_of_wisdom/pkg/lifecycle"
	"github.com/ognick/word_of_wisdom/pkg/logger"
	"github.com/ognick/word_of_wisdom/pkg/shutdown"
	"github.com/ognick/word_of_wisdom/pkg/tcp"
)

type App struct {
	log logger.Logger
	lc  lifecycle.Lifecycle
}

// All modules that need to be created and their Run method called.
type runnableModules struct {
	httpServer *http.Server
	tcpServer  *tcp.Server
}

func NewApp(
	log logger.Logger,
	lc lifecycle.Lifecycle,
	// side effects for creation components via wire
	_ runnableModules,
) *App {
	return &App{
		log: log,
		lc:  lc,
	}
}

func (app *App) Run() {
	runner, gracefulCtx := shutdown.CreateRunnerWithGracefulContext()
	app.lc.RunAllComponents(runner, gracefulCtx)
	app.log.Infof("Application started")
	// Awaiting graceful shutdown
	if err := runner.Wait(); err != nil {
		log.Fatalf("%v", err)
	}
	app.log.Infof("Application gracefully finished")
}
