package wisdom

import (
	"github.com/ognick/word_of_wisdom/internal/server/internal/domain/models"
	"github.com/ognick/word_of_wisdom/internal/server/internal/domain/types/repositories"
)

type Usecase struct {
	wisdomRepo repositories.Wisdom
}

func NewUsecase(wisdomRepo repositories.Wisdom) *Usecase {
	return &Usecase{wisdomRepo: wisdomRepo}
}

func (u *Usecase) GetWisdom() models.Wisdom {
	return u.wisdomRepo.GetWisdom()
}
