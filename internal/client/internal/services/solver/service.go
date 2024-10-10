package solver

type ProofOfWorkSolver interface {
	Solve(challenge []byte) ([]byte, error)
}

type Service struct {
	pow ProofOfWorkSolver
}

func NewService(pow ProofOfWorkSolver) *Service {
	return &Service{
		pow: pow,
	}
}

func (s *Service) Solve(challenge []byte) ([]byte, error) {
	return s.pow.Solve(challenge)
}
