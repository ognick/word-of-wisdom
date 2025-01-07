package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ognick/word_of_wisdom/pkg/logger"

	"github.com/ognick/word_of_wisdom/internal/server/internal/domain/types/usecases"
	"github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/api/http/v1/dto"
)

type Handler struct {
	log              logger.Logger
	challengeUsecase usecases.Challenge
	wisdomUsecase    usecases.Wisdom
}

func NewHandler(
	log logger.Logger,
	challengeUsecase usecases.Challenge,
	wisdomUsecase usecases.Wisdom,
) *Handler {
	return &Handler{
		log:              log,
		challengeUsecase: challengeUsecase,
		wisdomUsecase:    wisdomUsecase,
	}
}

func (h *Handler) Register(router gin.IRouter) {
	v1 := router.Group("/v1", proofOfWorkLimiter(h.log, h.challengeUsecase))
	{
		v1.GET("/wisdom", func(c *gin.Context) {
			response, err := json.Marshal(dto.NewWisdom(h.wisdomUsecase.GetWisdom()))
			if err != nil {
				c.String(http.StatusInternalServerError, "failed to marshal response")
				return
			}
			c.String(http.StatusOK, string(response))
		})
	}
}
