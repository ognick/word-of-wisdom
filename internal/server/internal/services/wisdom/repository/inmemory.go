package repository

import (
	"math/rand"

	"github.com/ognick/word_of_wisdom/internal/server/internal/domain/models"
)

type Repository struct {
	words []string
}

func NewRepository() *Repository {
	return &Repository{
		words: []string{
			"Don't communicate by sharing memory, share memory by communicating",
			"Concurrency is not parallelism",
			"Channels orchestrate; mutexes serialize",
			"The bigger the interface, the weaker the abstraction",
			"Make the zero value useful",
			"interface{} says nothing",
			"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite",
			"A little copying is better than a little dependency",
			"Syscall must always be guarded with build tags",
			"Cgo must always be guarded with build tags",
			"Errors are values",
			"Don't just check errors, handle them gracefully",
			"Design the architecture, name the components, document the details",
			"Documentation is for users",
		},
	}
}

func (repo *Repository) GetWisdom() models.Wisdom {
	size := len(repo.words)
	pos := rand.Intn(size)
	return models.Wisdom{
		Content: repo.words[pos],
	}
}
