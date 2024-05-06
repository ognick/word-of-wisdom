package app

import (
	"context"
	"errors"
	"fmt"
	net "net/http"
	"time"

	httpV1 "word_of_wisdom/internal/client/internal/api/http/v1"
	tcpV1 "word_of_wisdom/internal/client/internal/api/tcp/v1"
	"word_of_wisdom/internal/client/internal/service"
	"word_of_wisdom/internal/common/config"
	"word_of_wisdom/internal/common/constants"
	"word_of_wisdom/pkg/http"
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

	// TCP Handler
	tcpHandler := tcpV1.NewHandler(solverService)

	// HTTP Handler
	httpHandler := httpV1.NewHandler(solverService)

	// TCP Client
	tcpClient := tcp.NewClient(cfg.TCPAddress, tcpHandler.Handle)

	// HTTP Client
	var httpClient *http.Client
	httpClient = http.NewClient(
		cfg.HTTPAddress+"/v1/wisdom",
		"GET",
		func(ctx context.Context, status int, header net.Header, body []byte) error {
			if status != net.StatusOK {
				return fmt.Errorf("unknown status:%d body:%s", status, body)
			}

			if header.Get(constants.ChallengeHeader) == "" {
				log.Infof("wisdom:%s", body)
				return nil
			}

			if err := httpHandler.HandleChallenge(header); err != nil {
				return err
			}

			return httpClient.Request(ctx, header, nil)
		},
	)

	runner, ctx := shutdown.CreateRunnerWithGracefulContext()
	// Running tcp
	runner.Go(func() error {
		for {
			if err := tcpClient.Connect(ctx); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				log.Errorf("failed to connect: %v", err)
				<-time.After(cfg.ChallengeTimeout)
			}
		}
	})

	// Running http
	runner.Go(func() error {
		for {
			if err := httpClient.Request(ctx, nil, nil); err != nil {
				if err != nil {
					log.Errorf("failed to connect: %v", err)
				}
				select {
				case <-time.After(cfg.ChallengeTimeout):
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	})

	// Awaiting graceful shutdown
	if err := runner.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("%v", err)
	}

	log.Infof("Client gracefully finished")
}
