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

type Modules struct {
	httpServer *http.Server
	tcpServer  *tcp.Server
}

func NewApp(
	log logger.Logger,
	lc lifecycle.Lifecycle,
	// side effects for creation components via wire
	_ Modules,
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
