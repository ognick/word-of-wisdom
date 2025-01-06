package app

import (
	"log"

	"github.com/ognick/word_of_wisdom/pkg/http"
	"github.com/ognick/word_of_wisdom/pkg/logger"
	"github.com/ognick/word_of_wisdom/pkg/shutdown"
	"github.com/ognick/word_of_wisdom/pkg/tcp"
)

type App struct {
	log      logger.Logger
	tcpSrv   *tcp.Server
	httpServ *http.Server
}

func NewApp(log logger.Logger, tcpSrv *tcp.Server, httpServ *http.Server) *App {
	return &App{
		log:      log,
		tcpSrv:   tcpSrv,
		httpServ: httpServ,
	}
}

func (app *App) Run() {
	runner, gracefulCtx := shutdown.CreateRunnerWithGracefulContext()
	runner.Go(func() error {
		return app.tcpSrv.Run(gracefulCtx)
	})
	runner.Go(func() error {
		return app.httpServ.Run(gracefulCtx)
	})

	app.log.Infof("Server started")

	// Awaiting graceful shutdown
	if err := runner.Wait(); err != nil {
		log.Fatalf("%v", err)
	}

	app.log.Infof("Server gracefully finished")
}
