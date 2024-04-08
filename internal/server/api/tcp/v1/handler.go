package v1

import (
	"context"
	"errors"
	"net"
	"time"

	"word_of_wisdom/pkg/logger"
)

const bufSize = 64

type WisdomService interface {
	GetWisdom() string
}

type ChallengeService interface {
	GenerateChallenge() ([]byte, error)
	ValidateSolution(challenge, solution []byte) bool
}

type Handler struct {
	challengeService ChallengeService
	wisdomService    WisdomService
	timeout          time.Duration
	log              logger.Logger
}

func NewHandler(
	challengeService ChallengeService,
	wisdomService WisdomService,
	timeout time.Duration,
) *Handler {
	return &Handler{
		challengeService: challengeService,
		wisdomService:    wisdomService,
		timeout:          timeout,
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

	challenge, err := h.challengeService.GenerateChallenge()
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

	if !h.challengeService.ValidateSolution(challenge, solution[:n]) {
		h.log.Debugf("client couldn't proof of work")
	}

	wisdom := h.wisdomService.GetWisdom()
	_, err = conn.Write([]byte(wisdom))
	if err != nil {
		h.log.Errorf("failed to send wisdom: %v", err)
		return
	}

	h.log.Debugf("wisdom: %s", wisdom)
}
