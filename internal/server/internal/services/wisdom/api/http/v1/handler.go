package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ognick/word_of_wisdom/internal/server/internal/domain/interfaces/usecases"
	"github.com/ognick/word_of_wisdom/pkg/logger"
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

// getWisdom godoc
// @Summary      Get wisdom
// @Description  Get wisdom
// @Tags         V1
// @Accept       json
// @Produce      json
// @Param X-Challenge header string false "Header containing the generated challenge (provided only on the first request)."
// @Param X-Solution header string false "Header containing the solution to the challenge for validation (for subsequent requests)."
// @Success      200  {object}  dto.Wisdom
// @Failure      500  {string}  http.StatusInternalServerError
// @Failure      400  {string}  http.StatusBadRequest
// @Router       /v1/wisdom [get]
func (h *Handler) getWisdom(c *gin.Context) {
	response, err := json.Marshal(NewWisdom(h.wisdomUsecase.GetWisdom()))
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to marshal response")
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) Register(router gin.IRouter) {
	v1 := router.Group("/v1", proofOfWorkLimiter(h.log, h.challengeUsecase))
	{
		v1.GET("/wisdom", h.getWisdom)
	}
}
