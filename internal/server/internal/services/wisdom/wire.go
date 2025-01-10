package wisdom

import (
	"github.com/google/wire"

	"github.com/ognick/word_of_wisdom/internal/server/internal/domain/interfaces/usecases"
	httpv1 "github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/api/http/v1"
	tcpv1 "github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/api/tcp/v1"
	repo "github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/repository"
	usecase "github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/usecase"
)

func ProvideRepo() usecase.Repo {
	return repo.NewRepository()
}

func ProvideUsecase(repo usecase.Repo) usecases.Wisdom {
	return usecase.NewUsecase(repo)
}

var Init = wire.NewSet(
	ProvideRepo,
	ProvideUsecase,
	tcpv1.Init,
	httpv1.Init,
)
