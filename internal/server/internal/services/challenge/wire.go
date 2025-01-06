package challenge

import (
	"github.com/google/wire"

	"github.com/ognick/word_of_wisdom/internal/server/internal/domain/types/usecases"
	usecase "github.com/ognick/word_of_wisdom/internal/server/internal/services/challenge/usecase"
)

func ProvideUsecase(pow usecase.ProofOfWorkGenerator) usecases.Challenge {
	return usecase.NewUsecase(pow)
}

var Init = wire.NewSet(
	ProvideUsecase,
)
