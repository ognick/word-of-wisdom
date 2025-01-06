package usecases

import "github.com/ognick/word_of_wisdom/internal/server/internal/domain/models"

type Wisdom interface {
	GetWisdom() models.Wisdom
}
