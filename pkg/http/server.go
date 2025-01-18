package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ognick/word_of_wisdom/pkg/lifecycle"
	"github.com/ognick/word_of_wisdom/pkg/logger"
)

type Addr string

type Server struct {
	log logger.Logger
	srv *http.Server
}

func NewServer(lc lifecycle.Lifecycle, log logger.Logger, addr Addr, handler http.Handler) *Server {
	return lifecycle.RegisterComponent(lc,
		&Server{
			log: log,
			srv: &http.Server{
				Addr:    string(addr),
				Handler: handler,
			},
		})
}

func (s *Server) waitForReady(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Millisecond):
			conn, err := net.Dial("tcp", s.srv.Addr)
			if err != nil {
				return err
			}
			return conn.Close()
		}
	}
}

func (s *Server) Run(ctx context.Context, readinessProbe chan error) error {
	done := make(chan error)
	go func() {
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			done <- fmt.Errorf("error occurred while running http server: %w", err)
		}
		close(done)
	}()

	go func() {
		err := s.waitForReady(ctx)
		if err == nil {
			s.log.Infof("HTTP server was started on %s", s.srv.Addr)
		}
		readinessProbe <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			s.log.Errorf("failed to listen: %v", err)
		}
	case <-ctx.Done():
		if err := s.srv.Shutdown(ctx); err != nil {
			s.log.Errorf("failed to close listener: %v", err)
		}
	}

	s.log.Infof("HTTP server was stopped")

	return nil
}
