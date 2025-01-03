package wisdom

import "github.com/google/wire"

type WisdomRepo interface {
	GetWisdom() string
}

type Usecase struct {
	wisdomRepo WisdomRepo
}

func NewUsecase(wisdomRepo WisdomRepo) *Usecase {
	return &Usecase{wisdomRepo: wisdomRepo}
}

func (u *Usecase) GetWisdom() string {
	return u.wisdomRepo.GetWisdom()
}

var Set = wire.NewSet(NewUsecase)
