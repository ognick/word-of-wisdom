package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
}

func NewHandler(
	challengeService ChallengeService,
	wisdomService WisdomService,
) *Handler {
	return &Handler{
		challengeService: challengeService,
		wisdomService:    wisdomService,
	}
}

func (h *Handler) Init() *gin.Engine {
	// Init gin handler
	router := gin.Default()

	router.Use(
		gin.Recovery(),
	)

	// Init router
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	v1 := router.Group("/v1", ProofOfWorkLimiter(h.challengeService))
	{
		v1.GET("/wisdom", func(c *gin.Context) {
			c.String(http.StatusOK, h.wisdomService.GetWisdom())
		})
	}
}
