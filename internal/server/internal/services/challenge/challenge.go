package challenge

type ProofOfWorkGenerator interface {
	Generate() ([]byte, error)
	Validate([]byte, []byte) bool
}

type Service struct {
	pow ProofOfWorkGenerator
}

func NewService(pow ProofOfWorkGenerator) *Service {
	return &Service{
		pow: pow,
	}
}

func (s *Service) GenerateChallenge() ([]byte, error) {
	return s.pow.Generate()
}

func (s *Service) ValidateSolution(challenge, solution []byte) bool {
	return s.pow.Validate(challenge, solution)
}
