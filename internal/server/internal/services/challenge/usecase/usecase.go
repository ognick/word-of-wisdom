package challenge

import "github.com/google/wire"

type ProofOfWorkGenerator interface {
	Generate() ([]byte, error)
	Validate([]byte, []byte) bool
}

type Usecase struct {
	pow ProofOfWorkGenerator
}

func NewUsecase(pow ProofOfWorkGenerator) *Usecase {
	return &Usecase{
		pow: pow,
	}
}

func (u *Usecase) GenerateChallenge() ([]byte, error) {
	return u.pow.Generate()
}

func (u *Usecase) ValidateSolution(challenge, solution []byte) bool {
	return u.pow.Validate(challenge, solution)
}

var Set = wire.NewSet(NewUsecase)
