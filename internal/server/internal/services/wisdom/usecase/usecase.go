package wisdom

import (
	"github.com/ognick/word_of_wisdom/internal/server/internal/domain/models"
)

type Repo interface {
	GetWisdom() models.Wisdom
}

type Usecase struct {
	wisdomRepo Repo
}

func NewUsecase(wisdomRepo Repo) *Usecase {
	return &Usecase{wisdomRepo: wisdomRepo}
}

func (u *Usecase) GetWisdom() models.Wisdom {
	return u.wisdomRepo.GetWisdom()
}
