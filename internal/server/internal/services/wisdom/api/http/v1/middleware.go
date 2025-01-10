package v1

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/ognick/word_of_wisdom/internal/common/constants"
	"github.com/ognick/word_of_wisdom/internal/server/internal/domain/interfaces/usecases"
	"github.com/ognick/word_of_wisdom/pkg/logger"
)

const (
	cacheSize = 1000
	cacheTTL  = 1 * time.Second
)

func proofOfWorkLimiter(log logger.Logger, challengeUsecase usecases.Challenge) gin.HandlerFunc {
	cache := expirable.NewLRU[string, []byte](cacheSize, nil, cacheTTL)

	generateChallenge := func(c *gin.Context, addr string) {
		challenge, err := challengeUsecase.GenerateChallenge()
		if err != nil {
			log.Errorf("faile to generate chellange: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		cache.Add(addr, challenge)

		base64Challenge := base64.StdEncoding.EncodeToString(challenge)
		c.Header(constants.ChallengeHeader, base64Challenge)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	return func(c *gin.Context) {
		addr := c.Request.RemoteAddr
		base64Solution := c.GetHeader(constants.SolutionHeader)
		if base64Solution == "" {
			generateChallenge(c, addr)
			return
		}

		challenge, ok := cache.Get(addr)
		if !ok {
			generateChallenge(c, addr)
			return
		}

		solution, err := base64.StdEncoding.DecodeString(base64Solution)
		if err != nil {
			generateChallenge(c, addr)
			return
		}

		if !challengeUsecase.ValidateSolution(challenge, solution) {
			generateChallenge(c, addr)
			return
		}

		c.Next()
	}
}
