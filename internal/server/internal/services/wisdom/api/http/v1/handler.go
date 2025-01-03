package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"word_of_wisdom/internal/server/internal/domain/types/usecases"
)

type Handler struct {
	challengeUsecase usecases.Challenge
	wisdomUsecase    usecases.Wisdom
}

func NewHandler(
	challengeUsecase usecases.Challenge,
	wisdomUsecase usecases.Wisdom,
) *Handler {
	return &Handler{
		challengeUsecase: challengeUsecase,
		wisdomUsecase:    wisdomUsecase,
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
