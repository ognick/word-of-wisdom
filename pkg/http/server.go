package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"word_of_wisdom/pkg/logger"
)

type Server struct {
	log logger.Logger
	srv *http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		log: logger.NewLogger(),
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	done := make(chan error)
	go func() {
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			done <- fmt.Errorf("error occurred while running http server: %w", err)
		}
		close(done)
	}()

	s.log.Infof("HTTP server was started on %s", s.srv.Addr)

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
