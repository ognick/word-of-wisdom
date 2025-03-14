package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/ognick/word_of_wisdom/pkg/lifecycle"
	"github.com/ognick/word_of_wisdom/pkg/logger"
)

var opError = &net.OpError{}

type Addr string

type Server struct {
	addr    string
	handler func(conn net.Conn)
	log     logger.Logger
}

func NewServer(
	lc lifecycle.Lifecycle,
	log logger.Logger,
	addr Addr,
	handler func(conn net.Conn),
) *Server {
	return lifecycle.RegisterComponent(lc,
		&Server{
			addr:    string(addr),
			handler: handler,
			log:     log,
		})
}

func (srv *Server) Run(ctx context.Context, readinessProbe chan error) error {
	listener, err := net.Listen("tcp", srv.addr)
	if err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	go func() {
		<-ctx.Done()
		if err := listener.Close(); err != nil {
			srv.log.Errorf("failed to close listener: %v", err)
		}
	}()

	close(readinessProbe)
	srv.log.Infof("TCP server was starded on %v", listener.Addr())

	connections := make(chan net.Conn)
	go func() {
		defer close(connections)
		for {
			conn, err := listener.Accept()
			if err != nil {
				if !errors.As(err, &opError) {
					srv.log.Errorf("failed to accept: %v", err)
				}

				return
			}
			connections <- conn
		}
	}()

	for conn := range connections {
		go srv.handler(conn)
	}

	srv.log.Infof("TCP server was stopped")

	return nil
}
