package v1

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/google/wire"
	"word_of_wisdom/pkg/logger"
)

const bufSize = 64

type ChallengeTimeout time.Duration

type WisdomUsecase interface {
	GetWisdom() string
}

type ChallengeUsecase interface {
	GenerateChallenge() ([]byte, error)
	ValidateSolution(challenge, solution []byte) bool
}

type Handler struct {
	challengeUsecase ChallengeUsecase
	wisdomUsecase    WisdomUsecase
	timeout          time.Duration
	log              logger.Logger
}

func NewHandler(
	challengeUsecase ChallengeUsecase,
	wisdomUsecase WisdomUsecase,
	timeout ChallengeTimeout,
) *Handler {
	return &Handler{
		challengeUsecase: challengeUsecase,
		wisdomUsecase:    wisdomUsecase,
		timeout:          time.Duration(timeout),
		log:              logger.NewLogger(),
	}
}

func (h *Handler) Handle(conn net.Conn) {
	h.log.Debugf("handling")
	defer func() {
		if err := conn.Close(); err != nil {
			h.log.Errorf("connection closed, err: %v", err)
			return
		}
		h.log.Debugf("connection closed")
	}()

	challenge, err := h.challengeUsecase.GenerateChallenge()
	if err != nil {
		h.log.Errorf("failed to generate challange: %v", err)
		return
	}
	h.log.Debugf("challenge was created: %d bytes", len(challenge))
	if _, err := conn.Write(challenge); err != nil {
		h.log.Errorf("failed to send challange: %v", err)
		return
	}

	if err := conn.SetReadDeadline(time.Now().Add(h.timeout)); err != nil {
		h.log.Errorf("failed to set read deadline: %v", err)
		return
	}

	h.log.Debugf("waiting solving")
	solution := make([]byte, bufSize)
	n, err := conn.Read(solution)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.log.Debugf("challenge didn't submit")
			return
		}
		h.log.Errorf("error reading response: %v", err)
		return
	}

	if !h.challengeUsecase.ValidateSolution(challenge, solution[:n]) {
		h.log.Debugf("client couldn't proof of work")
		return
	}

	wisdom := h.wisdomUsecase.GetWisdom()
	_, err = conn.Write([]byte(wisdom))
	if err != nil {
		h.log.Errorf("failed to send wisdom: %v", err)
		return
	}

	h.log.Debugf("wisdom: %s", wisdom)
}

func ProvideTCPHandle(handler *Handler) func(net.Conn) {
	return handler.Handle
}

var Set = wire.NewSet(
	NewHandler,
	ProvideTCPHandle,
)
