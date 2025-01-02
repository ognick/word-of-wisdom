package challenge

type proofOfWorkGenerator interface {
	Generate() ([]byte, error)
	Validate([]byte, []byte) bool
}

type Usecase struct {
	pow proofOfWorkGenerator
}

func NewUsecase(pow proofOfWorkGenerator) *Usecase {
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
