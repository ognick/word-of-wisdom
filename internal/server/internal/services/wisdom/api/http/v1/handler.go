package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
}

func NewHandler(
	challengeService ChallengeUsecase,
	wisdomService WisdomUsecase,
) *Handler {
	return &Handler{
		challengeUsecase: challengeService,
		wisdomUsecase:    wisdomService,
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
	v1 := router.Group("/v1", ProofOfWorkLimiter(h.challengeUsecase))
	{
		v1.GET("/wisdom", func(c *gin.Context) {
			c.String(http.StatusOK, h.wisdomUsecase.GetWisdom())
		})
	}
}

func ProvideHTTPEngine(handler *Handler) *gin.Engine {
	return handler.Init()
}
