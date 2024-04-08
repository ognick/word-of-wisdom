package v1

import (
	"fmt"
	"net"

	"word_of_wisdom/pkg/logger"
)

const (
	bufSize = 64
)

// SolverService defines the interface for the task-solving service.
type SolverService interface {
	Solve(challenge []byte) ([]byte, error)
}

// Handler handles connections and solves tasks.
type Handler struct {
	solverService SolverService
	log           logger.Logger
}

// NewHandler creates a new instance of a handler with the given task-solving service.
func NewHandler(
	solverService SolverService,
) *Handler {
	return &Handler{
		solverService: solverService,
		log:           logger.NewLogger(),
	}
}

// Handle handles an incoming connection.
func (h *Handler) Handle(conn net.Conn) error {
	// Reading challenge from the connection.
	buf := make([]byte, bufSize)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to receive challenge: %w", err)
	}

	// Solving the challenge.
	solution, err := h.solverService.Solve(buf[:n])
	if err != nil {
		return fmt.Errorf("failed to solve challenge: %w", err)
	}

	// Sending solution back into the connection.
	_, err = conn.Write([]byte(solution))
	if err != nil {
		return fmt.Errorf("failed to send solution: %w", err)
	}

	// Reading wisdom from the connection.
	n, err = conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to receive wisdom feedback: %w", err)
	}

	h.log.Infof("Wisdom: '%s'", string(buf[:n]))
	return nil
}
