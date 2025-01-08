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
	s := &Server{
		log: log,
		srv: &http.Server{
			Addr:    string(addr),
			Handler: handler,
		},
	}
	lc.Register(s)
	return s
}

func (s *Server) waitForReady(ctx context.Context) bool {
	for {
		select {
		case <-ctx.Done():
			return false
		case <-time.After(1 * time.Millisecond):
			conn, err := net.Dial("tcp", s.srv.Addr)
			if err == nil {
				_ = conn.Close()
				return true
			}
		}
	}
}

func (s *Server) Run(ctx context.Context, ready chan struct{}) error {
	done := make(chan error)
	go func() {
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			done <- fmt.Errorf("error occurred while running http server: %w", err)
		}
		close(done)
	}()

	if s.waitForReady(ctx) {
		s.log.Infof("HTTP server was started on %s", s.srv.Addr)
		close(ready)
	}

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
