package v1

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/ognick/word_of_wisdom/internal/common/constants"
	"github.com/ognick/word_of_wisdom/pkg/logger"
)

type SolverService interface {
	Solve(challenge []byte) ([]byte, error)
}

type Handler struct {
	solverService SolverService
	log           logger.Logger
}

func NewHandler(
	solverService SolverService,
) *Handler {
	return &Handler{
		solverService: solverService,
		log:           logger.NewLogger(),
	}
}

func (h *Handler) HandleChallenge(header http.Header) error {
	base64Challenge := header.Get(constants.ChallengeHeader)
	challenge, err := base64.StdEncoding.DecodeString(base64Challenge)
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}
	solution, err := h.solverService.Solve(challenge)
	if err != nil {
		return fmt.Errorf("failed to solve: %w", err)
	}
	base64Solution := base64.StdEncoding.EncodeToString(solution)

	header.Set(constants.SolutionHeader, base64Solution)
	return nil
}
